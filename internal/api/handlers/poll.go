package handlers

import (
	"net/http"

	"mattermost-vote-bot/internal/service"
	"mattermost-vote-bot/pkg/logger"

	"github.com/gin-gonic/gin"
)

type PollHandler struct {
	service service.PollService
	log     logger.Logger
}

func NewPollHandler(s service.PollService, log logger.Logger) *PollHandler {
	return &PollHandler{
		service: s,
		log:     log,
	}
}

// CreatePollRequest - DTO для создания опроса
type CreatePollRequest struct {
	Question  string            `json:"question" binding:"required"`
	Creator   string            `json:"creator" binding:"required"`
	ChannelID string            `json:"channel_id" binding:"required"`
	Options   map[string]string `json:"options" binding:"required,min=2"`
}

// Create - создание нового опроса
func (h *PollHandler) Create(c *gin.Context) {
	var req CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	poll, err := h.service.Create(c.Request.Context(), &service.CreatePollRequest{
		Question:  req.Question,
		Creator:   req.Creator,
		ChannelID: req.ChannelID,
		Options:   req.Options,
	})
	if err != nil {
		h.log.WithError(err).Error("Failed to create poll")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, poll)
}

// GetPoll - получение опроса по ID
func (h *PollHandler) GetPoll(c *gin.Context) {
	pollID := c.Param("id")
	userID := c.Query("user_id")

	poll, err := h.service.GetPoll(c.Request.Context(), pollID, userID)
	if err != nil {
		h.log.WithError(err).WithField("poll_id", pollID).Error("Failed to get poll")
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

// EndPoll - завершение опроса
func (h *PollHandler) EndPoll(c *gin.Context) {
	pollID := c.Param("id")
	userID := c.Query("user_id") // В реальном приложении из токена

	if err := h.service.EndPoll(c.Request.Context(), pollID, userID); err != nil {
		h.log.WithError(err).WithField("poll_id", pollID).Error("Failed to end poll")
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
