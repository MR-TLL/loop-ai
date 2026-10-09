package handlers

import (
	"net/http"
	"backend/repository"
	"github.com/gin-gonic/gin"
)

type ClinicHandler struct {
	repo repository.ClinicRepository
}

func NewClinicHandler(repo repository.ClinicRepository) *ClinicHandler {
	return &ClinicHandler{repo: repo}
}

func (h *ClinicHandler) ListClinics(c *gin.Context) {
	clinics, err := h.repo.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clinics)
}

func (h *ClinicHandler) GetClinic(c *gin.Context) {
	id := c.Param("clinicId")
	clinic, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Clinic not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clinic)
}
