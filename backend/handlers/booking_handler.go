package handlers

import (
	"net/http"
	"backend/controller"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	ctrl *controller.BookingController
}

func NewBookingHandler(ctrl *controller.BookingController) *BookingHandler {
	return &BookingHandler{ctrl: ctrl}
}

type holdReq struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *BookingHandler) HoldSlot(c *gin.Context) {
	slotID := c.Param("slotId")
	var req holdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.ctrl.HoldSlot(c.Request.Context(), req.UserID, slotID)
	if err != nil {
		if err == controller.ErrSlotNotAvailable {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

type bookReq struct {
	UserID    string `json:"user_id" binding:"required"`
	PaymentID string `json:"payment_id" binding:"required"`
}

func (h *BookingHandler) BookSlot(c *gin.Context) {
	slotID := c.Param("slotId")
	var req bookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.ctrl.BookSlot(c.Request.Context(), req.UserID, slotID, req.PaymentID)
	if err != nil {
		switch err {
		case controller.ErrNotYourReservation:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case controller.ErrInvalidSlotState:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case controller.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}
