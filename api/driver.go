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

type CreateDriverRequest struct {
	Password             string  `json:"password" binding:"required,min=6"`
	FullName             string  `json:"full_name" binding:"required"`
	DrivingLicense       string  `json:"driving_license" binding:"required"`
	Mobile               string  `json:"mobile" binding:"required"`
	CarID                int64   `json:"car_id" binding:"required"`
	CarType              string  `json:"car_type" binding:"required"`
	CarImage             string  `json:"car_image" binding:"required"`
	OnlineStatus         bool    `json:"online_status" binding:"required"`
	Rating               float64 `json:"rating" binding:"required"`
	ProfileStatus        int32   `json:"profile_status" binding:"required"`
	SubscriptionStatus   bool    `json:"subscription_status" binding:"required"`
	SubscriptionPackage  string  `json:"subscription_package" binding:"required"`
	SubscriptionAmount   string  `json:"subscription_amount" binding:"required"`
	SubscriptionValidity int32   `json:"subscription_validity" binding:"required"`
}

type DriverResponse struct {
	FullName             string  `json:"full_name"`
	DrivingLicense       string  `json:"driving_license"`
	Mobile               string  `json:"mobile"`
	CarID                int64   `json:"car_id"`
	CarType              string  `json:"car_type"`
	CarImage             string  `json:"car_image"`
	OnlineStatus         bool    `json:"online_status"`
	Rating               float64 `json:"rating"`
	ProfileStatus        int32   `json:"profile_status"`
	SubscriptionStatus   bool    `json:"subscription_status"`
	SubscriptionPackage  string  `json:"subscription_package"`
	SubscriptionAmount   string  `json:"subscription_amount"`
	SubscriptionValidity int32   `json:"subscription_validity"`
}

func newDriverResponse(user db.Driver) DriverResponse {
	return DriverResponse{
		FullName:             user.FullName,
		DrivingLicense:       user.DrivingLicense,
		Rating:               user.Rating,
		Mobile:               user.Mobile,
		CarID:                user.CarID,
		CarType:              user.CarType,
		CarImage:             user.CarImage,
		OnlineStatus:         user.OnlineStatus,
		ProfileStatus:        user.ProfileStatus,
		SubscriptionStatus:   user.SubscriptionStatus,
		SubscriptionPackage:  user.SubscriptionPackage,
		SubscriptionAmount:   user.SubscriptionAmount,
		SubscriptionValidity: user.SubscriptionValidity,
	}
}

func (server *Server) createDriver(ctx *gin.Context) {
	var req CreateDriverRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	arg := db.CreateDriverParams{
		HashedPassword:       hashedPass,
		FullName:             req.FullName,
		DrivingLicense:       req.DrivingLicense,
		Rating:               req.Rating,
		Mobile:               req.Mobile,
		CarID:                req.CarID,
		CarType:              req.CarType,
		CarImage:             req.CarImage,
		OnlineStatus:         req.OnlineStatus,
		ProfileStatus:        req.ProfileStatus,
		SubscriptionStatus:   req.SubscriptionStatus,
		SubscriptionPackage:  req.SubscriptionPackage,
		SubscriptionAmount:   req.SubscriptionAmount,
		SubscriptionValidity: req.SubscriptionValidity,
	}

	driver, err := server.store.CreateDriver(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, finalResponse(FinalResponse{
					Status:  false,
					Message: "This driver is already registered"}))
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	driverResponse := newDriverResponse(driver)

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  false,
		Message: "Driver created successfully",
		Data:    driverResponse}))
}

type GetDriverRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getDriver(ctx *gin.Context) {
	var req GetDriverRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	driver, err := server.store.GetDriver(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
				Status:  false,
				Message: "Driver not found"}))
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
		Data:    newDriverResponse(driver)}))

}

type LoginDriverRequest struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginDriverResponse struct {
	AccessToken string         `json:"access_token"`
	User        DriverResponse `json:"user"`
}

func (server *Server) loginDriver(ctx *gin.Context) {
	var req LoginDriverRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	driver, err := server.store.GetDriverByMobile(ctx, req.Mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
				Status:  false,
				Message: "Driver not found with this mobile number"}))
			return
		}

		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	err = utils.CheckPassword(req.Password, driver.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, finalResponse(FinalResponse{
			Status:  false,
			Message: "Invalid password"}))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(driver.Mobile, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	rsp := LoginDriverResponse{
		AccessToken: accessToken,
		User:        newDriverResponse(driver),
	}

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  false,
		Message: "Login successful",
		Data:    rsp}))
}
