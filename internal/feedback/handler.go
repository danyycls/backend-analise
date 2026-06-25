package feedback

import (
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Usecase *SaveFeedbackUsecase
}

func (h *Handler) SaveFeedback(c *gin.Context) {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid feedback payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}
	if err := h.Usecase.Execute(c.Request.Context(), req.Name, req.Email, req.Message); err != nil {
		logger.Error("failed to save feedback", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível salvar feedback"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Feedback enviado com sucesso"})
}
