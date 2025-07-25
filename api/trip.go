package api

import (
	"database/sql"
	"net/http"

	db "github.com/emonoid/toribook.git/db/sqlc"
	"github.com/emonoid/toribook.git/helpers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CreateTripRequest struct {
	BookingID       string  `json:"booking_id" binding:"required"`
	TripStatus      string  `json:"trip_status" binding:"required"`
	PickupLocation  string  `json:"pickup_location" binding:"required"`
	PickupLat       string  `json:"pickup_lat" binding:"required"`
	PickupLong      string  `json:"pickup_long" binding:"required"`
	DropoffLocation string  `json:"dropoff_location" binding:"required"`
	DropoffLat      string  `json:"dropoff_lat" binding:"required"`
	DropoffLong     string  `json:"dropoff_long" binding:"required"`
	DriverID        *int64  `json:"driver_id"`
	DriverName      *string `json:"driver_name"`
	DriverMobile    *string `json:"driver_mobile"`
	CarID           *int64  `json:"car_id"`
	CarType         *string `json:"car_type"`
	CarImage        *string `json:"car_image"`
	Fare            *int32  `json:"fare"`
}

type TripResponse struct {
	BookingID       string  `json:"booking_id"`
	TripStatus      string  `json:"trip_status"`
	PickupLocation  string  `json:"pickup_location"`
	PickupLat       string  `json:"pickup_lat"`
	PickupLong      string  `json:"pickup_long"`
	DropoffLocation string  `json:"dropoff_location"`
	DropoffLat      string  `json:"dropoff_lat"`
	DropoffLong     string  `json:"dropoff_long"`
	DriverID        *int64  `json:"driver_id"`
	DriverName      *string `json:"driver_name"`
	DriverMobile    *string `json:"driver_mobile"`
	CarID           *int64  `json:"car_id"`
	CarType         *string `json:"car_type"`
	CarImage        *string `json:"car_image"`
	Fare            *int32  `json:"fare"`
}

func (server *Server) createTrip(ctx *gin.Context) {
	var req CreateTripRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	arg := db.CreateTripParams{
		BookingID:       req.BookingID,
		TripStatus:      req.TripStatus,
		PickupLocation:  req.PickupLocation,
		PickupLat:       req.PickupLat,
		PickupLong:      req.PickupLong,
		DropoffLocation: req.DropoffLocation,
		DropoffLat:      req.DropoffLat,
		DropoffLong:     req.DropoffLong,
		DriverID:        helpers.MakeNullInt64(req.DriverID),
		DriverName:      helpers.MakeNullString(req.DriverName),
		DriverMobile:    helpers.MakeNullString(req.DriverMobile),
		CarID:           helpers.MakeNullInt64(req.CarID),
		CarType:         helpers.MakeNullString(req.CarType),
		CarImage:        helpers.MakeNullString(req.CarImage),
		Fare:            helpers.MakeNullInt32(req.Fare),
	}

	trip, err := server.store.CreateTrip(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	response := TripResponse{
		BookingID:       trip.BookingID,
		TripStatus:      trip.TripStatus,
		PickupLocation:  trip.PickupLocation,
		PickupLat:       trip.PickupLat,
		PickupLong:      trip.PickupLong,
		DropoffLocation: trip.DropoffLocation,
		DropoffLat:      trip.DropoffLat,
		DropoffLong:     trip.DropoffLong,
		DriverID:        helpers.NullInt64ToPtr(trip.DriverID),
		DriverName:      helpers.NullStringToPtr(trip.DriverName),
		DriverMobile:    helpers.NullStringToPtr(trip.DriverMobile),
		CarID:           helpers.NullInt64ToPtr(trip.CarID),
		CarType:         helpers.NullStringToPtr(trip.CarType),
		CarImage:        helpers.NullStringToPtr(trip.CarImage),
		Fare:            helpers.NullInt32ToPtr(trip.Fare),
	}

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  false,
		Message: "Trip created successfully",
		Data:    response}))

	server.webSocketManager.Broadcast("trips", finalResponse(FinalResponse{
		Status:  false,
		Message: "Trip created successfully",
		Data:    response}))
}

type GetTripRequest struct {
	BookingID string `uri:"booking_id" binding:"required"`
}

func (server *Server) getTrip(ctx *gin.Context) {
	var req GetTripRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	trip, err := server.store.GetTripByBookingID(ctx, req.BookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
				Status:  false,
				Message: "Trip not found"}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	response := TripResponse{
		BookingID:       trip.BookingID,
		TripStatus:      trip.TripStatus,
		PickupLocation:  trip.PickupLocation,
		PickupLat:       trip.PickupLat,
		PickupLong:      trip.PickupLong,
		DropoffLocation: trip.DropoffLocation,
		DropoffLat:      trip.DropoffLat,
		DropoffLong:     trip.DropoffLong,
		DriverID:        helpers.NullInt64ToPtr(trip.DriverID),
		DriverName:      helpers.NullStringToPtr(trip.DriverName),
		DriverMobile:    helpers.NullStringToPtr(trip.DriverMobile),
		CarID:           helpers.NullInt64ToPtr(trip.CarID),
		CarType:         helpers.NullStringToPtr(trip.CarType),
		CarImage:        helpers.NullStringToPtr(trip.CarImage),
		Fare:            helpers.NullInt32ToPtr(trip.Fare),
	}

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  false,
		Message: "Success",
		Data:    response}))

}

type GetAllTripsRequest struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PerPage    int32 `form:"per_page" binding:"required"`
}

func (server *Server) getAllTrips(ctx *gin.Context) {
	var req GetAllTripsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error(),
			Data:    nil}))
		return
	}

	arg := db.ListTripsParams{
		Limit:  req.PerPage,
		Offset: (req.PageNumber - 1) * req.PerPage,
	}

	trips, err := server.store.ListTrips(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
				Status:  false,
				Message: "No trips found",
				Data:    []db.Trip{}}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error(),
			Data:    []db.Trip{}},
		))
		return
	}

	if trips == nil {
		ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
			Status:  false,
			Message: "No trips found",
			Data:    []db.Trip{}},
		))
		return
	}

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  true,
		Message: "Trips retrieved successfully",
		Data:    trips,
	}))

}

type UpdateTripStatusRequest struct {
	BookingID  string `json:"booking_id" binding:"required"`
	TripStatus string `json:"trip_status" binding:"required"`
}

func (server *Server) updateTripStatus(ctx *gin.Context) {
	var req UpdateTripStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	trip, err := server.store.UpdateTripStatus(ctx, db.UpdateTripStatusParams{
		BookingID:  req.BookingID,
		TripStatus: req.TripStatus,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
				Status:  false,
				Message: "Trip not found",
				Data:    nil}))
			return
		}

		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  true,
		Message: "Trip status updated successfully",
		Data:    trip,
	}))

	server.webSocketManager.Broadcast("trip_status: "+trip.BookingID, finalResponse(FinalResponse{
		Status:  true,
		Message: "Trip status updated",
		Data:    trip,
	}))
}

var tripUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (server *Server) tripWebSocket(ctx *gin.Context) {
	tokenString := ctx.Query("token")

	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, finalResponse(FinalResponse{
			Status:  false,
			Message: "Trip created successfully",
			Data:    "Missing token"}))
		return
	}

	_, err := server.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error()}))
		return
	}

	conn, err := tripUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	server.webSocketManager.AddClient("trips", conn)
	defer func() {
		server.webSocketManager.RemoveClient("trips", conn)
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

}

func (server *Server) tripStatusUpdateWebSocket(ctx *gin.Context) {
	tokenString := ctx.Query("token")
	bookingID := ctx.Query("booking_id")

	if tokenString == "" || bookingID == "" {
		ctx.JSON(http.StatusBadRequest, finalResponse(FinalResponse{
			Status:  false,
			Message: "Missing token or booking_id",
		}))
		return
	}

	_, err := server.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error(),
		}))
		return
	}

	conn, err := tripUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	server.webSocketManager.AddClient("trip_status: "+bookingID, conn)
	defer func() {
		server.webSocketManager.RemoveClient("trip_status: "+bookingID, conn)
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
