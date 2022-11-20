package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/registry"
	"github.com/lampctl/lampctl/ui"
	"github.com/nathan-osman/go-herald"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Server provides an HTTP interface for interacting with lamps.
type Server struct {
	server   http.Server
	herald   *herald.Herald
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
			herald:   herald.New(),
			logger:   log.With().Str("package", "server").Logger(),
			registry: cfg.Registry,
		}
	)

	// Serve the static files
	r.Use(static.Serve("/", ui.EmbedFileSystem{FileSystem: http.FS(ui.Content)}))

	// Define the API
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

	// Add the provider API routes
	api.GET("/providers", s.api_providers_GET)
	api.GET("/providers/:id", s.api_providers_id_GET)
	api.POST("/providers/:id/apply", s.api_providers_id_apply_POST)

	// Special route for websocket connections
	api.GET("/ws", s.api_ws_GET)

	// Add the API routes from each individual provider
	for _, p := range s.registry.Providers() {
		if err := p.Init(api); err != nil {
			return nil, err
		}
	}

	// Serve the static files on all other paths
	r.NoRoute(func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
		c.Abort()
	})

	// Set the handler for incoming socket messages and start the goroutine
	s.herald.MessageHandler = s.messageHandler
	s.herald.Start()

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

// Apply applies the provided changes to the specified provider.
func (s *Server) Apply(provider_id string, changes []*registry.Change) error {
	p, err := s.registry.GetProvider(provider_id)
	if err != nil {
		return err
	}
	return p.Apply(changes)
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Shutdown(context.Background())
}
