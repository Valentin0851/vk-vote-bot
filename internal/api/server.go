package api

import (
	"fmt"
	"net/http"
	"time"

	"mattermost-vote-bot/internal/service"
	"mattermost-vote-bot/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine      *gin.Engine
	pollService *service.PollService
	voteService *service.VoteService
	config      *Config
	log         logger.Logger
}

type Config struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func New(cfg *Config, pollService *service.PollService, voteService *service.VoteService, log logger.Logger) *Server {
	engine := gin.New()
	srv := &Server{
		engine:      engine,
		pollService: pollService,
		voteService: voteService,
		config:      cfg,
		log:         log,
	}
	srv.setupRoutes()
	return srv
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:      s.engine,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}
	return srv.ListenAndServe()
}
