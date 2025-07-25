package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
 
	"github.com/emonoid/toribook.git/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadkey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, finalResponse(FinalResponse{
				Status:  false,
				Message: err.Error()}))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, finalResponse(FinalResponse{
				Status:  false,
				Message: err.Error()}))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, finalResponse(FinalResponse{
				Status:  false,
				Message: err.Error()}))
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, finalResponse(FinalResponse{
				Status:  false,
				Message: err.Error()}))
			return
		}

		ctx.Set(authorizationPayloadkey, payload)

		ctx.Next()
	}
}
