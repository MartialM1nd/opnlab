package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/MartialM1nd/opnlab/internal/providers"
	"github.com/gin-gonic/gin"
)

// Server represents the opnlab API server.
type Server struct {
	router     *gin.Engine
	providers  map[string]providers.Provider
	providerMu sync.RWMutex
}

// New creates a new Server instance.
func New() *Server {
	s := &Server{
		providers: make(map[string]providers.Provider),
	}
	s.setupRouter()
	return s
}

// RegisterProvider adds a provider to the server.
func (s *Server) RegisterProvider(p providers.Provider) {
	s.providerMu.Lock()
	defer s.providerMu.Unlock()
	s.providers[p.Name()] = p
}

// setupRouter configures all HTTP routes.
func (s *Server) setupRouter() {
	s.router = gin.Default()

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   "",
		})
	})

	// API routes
	api := s.router.Group("/api")
	{
		api.GET("/providers", s.listProviders)
		api.GET("/providers/:name", s.getProvider)
		api.POST("/providers/:name/actions/:action", s.executeAction)
	}
}

// listProviders returns a list of all registered providers.
func (s *Server) listProviders(c *gin.Context) {
	s.providerMu.RLock()
	defier s.providerMu.RUnlock()
	names := make([]string, 0, len(s.providers))
	for name := range s.providers {
		names = append(names, name)
	}

	c.JSON(http.StatusOK, gin.H{
		"providers": names,
	})
}

// getProvider returns data from a specific provider.
func (s *Server) getProvider(c *gin.Context) {
	name := c.Param("name")

	s.providerMu.RLock()
	provider, exists := s.providers[name]
	s.providerMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("provider not found: %s", name),
		})
		return
	}

	data, err := provider.Collect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to collect data: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"data": data,
	})
}

// executeAction executes an action on a provider.
func (s *Server) executeAction(c *gin.Context) {
	name := c.Param("name")
	actionName := c.Param("action")

	s.providerMu.RLock()
	provider, exists := s.providers[name]
	s.providerMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("provider not found: %s", name),
		})
		return
	}

	actions := provider.Actions()
	action, exists := actions[actionName]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("action not found: %s", actionName),
		})
		return
	}

	// Parse request body as JSON params
	var params map[string]string
	if err := json.NewDecoder(c.Request.Body).Decode(¶ms); err != nil {
		params = make(map[string]string)
	}

	if err := action.Execute(params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("action failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("action %s executed", actionName),
	})
}

// Run starts the HTTP server.
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
