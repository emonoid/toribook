package api

import (
	"context"
	"encoding/json"
	"log"
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

func (server *Server) bidSubmitHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		server.bidSubmit(ctx, redisClient)
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

func (server *Server) getBidsHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		server.getBids(ctx, redisClient)
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

func (server *Server) BidWebSocketHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		server.bidWebSocket(ctx, redisClient)
	}
}

func (server *Server) bidWebSocket(ctx *gin.Context, redisClient *redis.Client) {
	bookingID := ctx.Param("booking_id")
	tokenString := ctx.Query("token")

	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	_, err := server.tokenMaker.VerifyToken(tokenString)
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

	channel := "bids_channel:" + bookingID
	server.webSocketManager.AddClient(channel, conn)
	defer func() {
		server.webSocketManager.RemoveClient(channel, conn)
		conn.Close()
	}()

	// Call listener only if not yet started
	server.redisLock.Lock()
	if !server.redisSubscribers[channel] {
		server.redisSubscribers[channel] = true
		server.StartBidChannelListener(redisClient, bookingID)
	}
	server.redisLock.Unlock()

	// Optional: keep connection alive with read pump
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (server *Server) StartBidChannelListener(redisClient *redis.Client, bookingID string) {
	channel := "bids_channel:" + bookingID

	go func() {
		sub := redisClient.Subscribe(context.Background(), channel)
		defer sub.Close()

		ch := sub.Channel()

		for msg := range ch {
			var bid Bid
			err := json.Unmarshal([]byte(msg.Payload), &bid)
			if err != nil {
				log.Println("Invalid bid payload:", err)
				continue
			}
			server.webSocketManager.Broadcast(channel, bid)
		}
	}()
}
