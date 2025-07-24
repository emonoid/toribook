package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/emonoid/toribook.git/db/sqlc"
	"github.com/emonoid/toribook.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreatePassengerRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type PassengerResponse struct {
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Rating   float64 `json:"rating"`
}

func newPassengerResponse(user db.Passenger) PassengerResponse {
	return PassengerResponse{
		FullName: user.FullName,
		Email:    user.Email,
		Rating:   user.Rating,
	}
}

func (server *Server) createPassenger(ctx *gin.Context) {
	var req CreatePassengerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if err.Error() == "Key: 'CreatePassengerRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag" {
			ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
				Status:  false,
				Message: "Invalid email format"}))
			return
		}
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	arg := db.CreatePassengerParams{
		HashedPassword: hashedPass,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	passenger, err := server.store.CreatePassenger(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, finalResponse(FinalResponse{
					Status:  false,
					Message: "This email is already registered"}))
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	passengerResponse := newPassengerResponse(passenger)

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  false,
		Message: "Passenger created successfully",
		Data:    passengerResponse}))
}

type GetPassengerRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getPassenger(ctx *gin.Context) {
	var req GetPassengerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	passenger, err := server.store.GetPassenger(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
				Status:  false,
				Message: err.Error()}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  false,
		Message: "Success",
		Data:    newPassengerResponse(passenger)}))

}

type LoginPassengerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginPassengerResponse struct {
	AccessToken string            `json:"access_token"`
	User        PassengerResponse `json:"user"`
}

func (server *Server) loginPassenger(ctx *gin.Context) {
	var req LoginPassengerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if err.Error() == "Key: 'LoginPassengerRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag" {
			ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
				Status:  false,
				Message: "Invalid email format"}))
			return
		}
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	passenger, err := server.store.GetPassengerByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
				Status:  false,
				Message: "User not found"}))
			return
		}

		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	err = utils.CheckPassword(req.Password, passenger.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, finalResponse(FinalResponse{
			Status:  false,
			Message: "Invalid password"}))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(passenger.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	rsp := LoginPassengerResponse{
		AccessToken: accessToken,
		User:        newPassengerResponse(passenger),
	}

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  false,
		Message: "Login successful",
		Data:    rsp}))
}
