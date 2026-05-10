package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"prompt-vcs/internal/models"
	storage "prompt-vcs/internal/storage/db"
)

type Server struct {
	echo *echo.Echo
	db   *storage.DB
}

func NewServer(dbPath string) (*Server, error) {
	db, err := storage.New(dbPath)
	if err != nil {
		return nil, err
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	s := &Server{
		echo: e,
		db:   db,
	}

	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {
	// Prompt CRUD
	s.echo.GET("/api/prompts", s.listPrompts)
	s.echo.POST("/api/prompts", s.createPrompt)
	s.echo.GET("/api/prompts/:id", s.getPrompt)
	s.echo.PUT("/api/prompts/:id", s.updatePrompt)
	s.echo.DELETE("/api/prompts/:id", s.deletePrompt)

	// Version control
	s.echo.GET("/api/prompts/:id/commits", s.getCommits)
	s.echo.POST("/api/prompts/:id/commits", s.commitPrompt)
	s.echo.GET("/api/prompts/:id/diff", s.diffPrompts)

	// Collections
	s.echo.GET("/api/collections", s.listCollections)
	s.echo.POST("/api/collections", s.createCollection)
}

func (s *Server) Start(address string) error {
	return s.echo.Start(address)
}

func (s *Server) listPrompts(c echo.Context) error {
	collection := c.QueryParam("collection")
	if collection == "" {
		collection = "default"
	}

	prompts, err := s.db.ListPrompts(collection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, prompts)
}

func (s *Server) createPrompt(c echo.Context) error {
	var prompt models.Prompt
	if err := c.Bind(&prompt); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	prompt.ID = generateID()
	prompt.CreatedAt = time.Now()
	prompt.UpdatedAt = time.Now()

	if err := s.db.CreatePrompt(&prompt); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, prompt)
}

func (s *Server) getPrompt(c echo.Context) error {
	id := c.Param("id")
	prompt, err := s.db.GetPrompt(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Prompt not found"})
	}
	return c.JSON(http.StatusOK, prompt)
}

func (s *Server) updatePrompt(c echo.Context) error {
	id := c.Param("id")
	var prompt models.Prompt
	if err := c.Bind(&prompt); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// In real implementation, merge with existing and create new version
	_ = id
	return c.JSON(http.StatusOK, prompt)
}

func (s *Server) deletePrompt(c echo.Context) error {
	id := c.Param("id")
	_ = id // Delete from DB
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) getCommits(c echo.Context) error {
	id := c.Param("id")
	commits, err := s.db.GetCommits(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, commits)
}

func (s *Server) commitPrompt(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		Message string `json:"message"`
		Author  string `json:"author"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Create commit
	hash := generateHash()
	if err := s.db.AddCommit(id, hash, req.Message, req.Author); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"hash": hash})
}

func (s *Server) diffPrompts(c echo.Context) error {
	// Implementation for diff
	return c.JSON(http.StatusOK, map[string]string{"diff": "coming soon"})
}

func (s *Server) listCollections(c echo.Context) error {
	// Implementation for listing collections
	return c.JSON(http.StatusOK, []models.Collection{})
}

func (s *Server) createCollection(c echo.Context) error {
	var collection models.Collection
	if err := c.Bind(&collection); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, collection)
}

func generateID() string {
	return "prompt-" + time.Now().Format("20060102150405")
}

func generateHash() string {
	// Simplified hash generation using timestamp
	return "commit-" + time.Now().Format("20060102150405")
}