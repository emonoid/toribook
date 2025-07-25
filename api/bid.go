package api

import (
	"encoding/json" 
	"net/http" 
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Bid struct {
	ID           string `json:"id"`
	BookingID    string `json:"booking_id"`
	BidAmount    int    `json:"bid_amount"`
	DriverID     int64  `json:"driver_id"`
	DriverName   string `json:"driver_name"`
	DriverRating int    `json:"driver_rating"`
	DriverMobile string `json:"driver_mobile"`
	CarID        int64  `json:"car_id"`
	CarType      string `json:"car_type"`
	CarImage     string `json:"car_image"`
}

func (s *Server) bidSubmitHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.bidSubmit(ctx, redisClient)
	}
}

func (server *Server) bidSubmit(ctx *gin.Context, redisClient *redis.Client) {
	var bid Bid
	if err := ctx.BindJSON(&bid); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid bid"})
		return
	}
	err := AddBid(redisClient, bid.BookingID, bid, ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to save bid"})
		return
	}
	_ = PublishBid(redisClient, bid.BookingID, bid, ctx)
	ctx.JSON(200, gin.H{"status": "bid placed"})
}

func (s *Server) getBidsHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.getBids(ctx, redisClient)
	}
}

func (server *Server) getBids(ctx *gin.Context, redisClient *redis.Client) {
	bookingID := ctx.Param("booking_id")
	bids, err := GetBids(redisClient, bookingID, ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch bids"})
		return
	}
	ctx.JSON(200, bids)
}

func AddBid(client *redis.Client, bookingID string, bid Bid, ctx *gin.Context) error {
	key := "bids:" + bookingID
	bidJSON, err := json.Marshal(bid)
	if err != nil {
		return err
	}
	err = client.RPush(ctx, key, bidJSON).Err()
	if err == nil {
		// client.Expire(ctx, key, time.Hour) // optional expiration
		client.Expire(ctx, key, 15*time.Minute)
	}
	return err
}

func GetBids(client *redis.Client, bookingID string, ctx *gin.Context) ([]Bid, error) {
	key := "bids:" + bookingID
	bidStrings, err := client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var bids []Bid
	for _, s := range bidStrings {
		var bid Bid
		if err := json.Unmarshal([]byte(s), &bid); err == nil {
			bids = append(bids, bid)
		}
	}
	return bids, nil
}

func PublishBid(client *redis.Client, bookingID string, bid Bid, ctx *gin.Context) error {
	channel := "bids_channel:" + bookingID
	msg, err := json.Marshal(bid)
	if err != nil {
		return err
	}
	return client.Publish(ctx, channel, msg).Err()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) BidWebSocketHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.bidWebSocket(ctx, redisClient)
	}
}

func (s *Server) bidWebSocket(ctx *gin.Context, redisClient *redis.Client) {

	bookingID := ctx.Param("booking_id")
	tokenString := ctx.Query("token")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}  

	_, err := s.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	sub := redisClient.Subscribe(ctx, "bids_channel:"+bookingID)
	ch := sub.Channel()

	for msg := range ch {
		conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
	}

}
