package handlers

import (
	"net/http"
	"strings"
	"backend/models"
	"backend/repository"
	"github.com/gin-gonic/gin"
)

type SlotHandler struct {
	repo repository.SlotRepository
}

func NewSlotHandler(repo repository.SlotRepository) *SlotHandler {
	return &SlotHandler{repo: repo}
}

func (h *SlotHandler) ListSlots(c *gin.Context) {
	clinicID := c.Param("clinicId")
	dateFilter := c.Query("date") // YYYY-MM-DD

	slots, err := h.repo.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if dateFilter != "" {
		filtered := []models.Slot{}
		for _, s := range slots {
			if strings.HasPrefix(s.StartTime.Format("2006-01-02"), dateFilter) {
				filtered = append(filtered, s)
			}
		}
		c.JSON(http.StatusOK, filtered)
		return
	}

	c.JSON(http.StatusOK, slots)
}
