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

	// //user routes
	// router.POST("/users", server.createUser)
	// router.GET("/users/:username", server.getUser)
	// router.POST("/users/login", server.loginUser)

	// protectedRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// // account routes
	// protectedRoutes.POST("/accounts", server.createAccount)
	// protectedRoutes.GET("/accounts/:id", server.getAccount)
	// protectedRoutes.GET("/accounts", server.getAllAccounts)
	// protectedRoutes.PUT("accounts/update", server.updateAccount)
	// protectedRoutes.DELETE("accounts/delete/:id", server.deleteAccount)

	// // transfer money routes
	// protectedRoutes.POST("/transfer", server.transferBalance)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
