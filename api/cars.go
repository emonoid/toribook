package api

import (
	"database/sql"
	"net/http"

	db "github.com/emonoid/toribook.git/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (server *Server) getAllCars(ctx *gin.Context) {
	cars, err := server.store.ListCars(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
				Status:  false,
				Message: "No cars found",
				Data:    []db.Car{}}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, finalResponse(FinalResponse{
			Status:  false,
			Message: err.Error(),
			Data:    []db.Car{}},
		))
		return
	}

	if cars == nil {
		ctx.JSON(http.StatusNotFound, finalResponse(FinalResponse{
			Status:  false,
			Message: "No cars found",
			Data:    []db.Car{}},
		))
		return
	}

	ctx.JSON(http.StatusOK, finalResponse(FinalResponse{
		Status:  true,
		Message: "Cars retrieved successfully",
		Data:    cars,
	}))

}
