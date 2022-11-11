package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/registry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Server provides an HTTP interface for interacting with lamps.
type Server struct {
	server   http.Server
	logger   zerolog.Logger
	registry *registry.Registry
}

func New(cfg *Config) (*Server, error) {

	// Switch to release mode if DEBUG is not set
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	var (
		r = gin.New()
		s = &Server{
			server: http.Server{
				Addr:    cfg.Addr,
				Handler: r,
			},
			logger:   log.With().Str("package", "server").Logger(),
			registry: cfg.Registry,
		}
	)

	api := r.Group("/api")

	// Attempt to handle panic() calls within API routes by converting them
	// into proper JSON responses
	api.Use(gin.CustomRecovery(func(c *gin.Context, i interface{}) {
		var message string
		switch v := i.(type) {
		case error:
			message = v.Error()
		case string:
			message = v
		default:
			message = "an unknown error has occurred"
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
	}))

	api.GET("/providers", s.api_providers_GET)
	api.GET("/providers/:id", s.api_providers_id_GET)

	// Start the goroutine that listens for incoming connections
	go func() {
		defer s.logger.Info().Msg("server stopped")
		s.logger.Info().Msg("server started")
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error().Msg(err.Error())
		}
	}()

	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Shutdown(context.Background())
}
