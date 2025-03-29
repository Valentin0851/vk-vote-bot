package handlers

import (
	"net/http"

	"mattermost-vote-bot/internal/service"
	"mattermost-vote-bot/pkg/logger"

	"github.com/gin-gonic/gin"
)

type VoteHandler struct {
	service service.VoteService
	log     logger.Logger
}

func NewVoteHandler(s service.VoteService, log logger.Logger) *VoteHandler {
	return &VoteHandler{
		service: s,
		log:     log,
	}
}

// VoteRequest - DTO для голосования
type VoteRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

// Vote - обработка голоса
func (h *VoteHandler) Vote(c *gin.Context) {
	pollID := c.Param("id")
	var req VoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("Invalid vote payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err := h.service.Vote(c.Request.Context(), &service.VoteRequest{
		PollID:   pollID,
		UserID:   req.UserID,
		OptionID: req.OptionID,
	})
	if err != nil {
		h.log.WithError(err).
			WithFields(map[string]interface{}{
				"poll_id":   pollID,
				"user_id":   req.UserID,
				"option_id": req.OptionID,
			}).
			Error("Failed to process vote")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetResults - получение результатов
func (h *VoteHandler) GetResults(c *gin.Context) {
	pollID := c.Param("id")
	userID := c.Query("user_id")

	results, err := h.service.GetResults(c.Request.Context(), pollID, userID)
	if err != nil {
		h.log.WithError(err).WithField("poll_id", pollID).Error("Failed to get results")
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	c.JSON(http.StatusOK, results)
}
