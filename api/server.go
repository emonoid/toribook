package api

import ( 
	"fmt"

	db "github.com/emonoid/toribook.git/db/sqlc"
	"github.com/emonoid/toribook.git/token"
	"github.com/emonoid/toribook.git/utils"
	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
	// "github.com/go-playground/validator/v10"
)

type Server struct {
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
}

func NewServer(config utils.Config, store *db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker([]byte(config.TokenSymmetricKey))
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	// Register custom validation if needed
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validateCurrency)
	// }

	server.setupRouters()

	return server, nil
}

func (server *Server) setupRouters() {
	router := gin.Default()
	protectedRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	redisClient := utils.NewRedisClient()

	apiVersion := "/api/v1/"
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, fmt.Sprintf("Welcome to Toribook %s", apiVersion))
	})

	// passenger routes
	router.POST(apiVersion+"passenger/registration", server.createPassenger)
	router.POST(apiVersion+"passenger/login", server.loginPassenger)
	protectedRoutes.GET(apiVersion+"passenger/:id", server.getPassenger)

	// driver routes
	router.POST(apiVersion+"driver/registration", server.createDriver)
	router.POST(apiVersion+"driver/login", server.loginDriver)
	protectedRoutes.GET(apiVersion+"driver/:id", server.getDriver)

	// trip routes
	protectedRoutes.POST(apiVersion+"trip/create", server.createTrip)
	protectedRoutes.GET(apiVersion+"trip/:id", server.getTrip)

	// bid routes
	protectedRoutes.POST(apiVersion+"bid/submit",server.bidSubmitHandler(redisClient))
	protectedRoutes.GET(apiVersion+"bids/:booking_id", server.getBidsHandler(redisClient))
	router.GET(apiVersion+"ws/bids/:booking_id", server.BidWebSocketHandler(redisClient))

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

type FinalResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func finalResponse(response FinalResponse) gin.H {
	return gin.H{"status": response.Status, "message": response.Message, "data": response.Data}
}
