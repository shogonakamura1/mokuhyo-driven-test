package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mokuhyo-driven-test/api/internal/handler"
	"github.com/mokuhyo-driven-test/api/internal/repository"
	"github.com/mokuhyo-driven-test/api/internal/service"
	"github.com/mokuhyo-driven-test/api/pkg/auth"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// JWKS client for JWT verification
	jwksURL := os.Getenv("SUPABASE_JWKS_URL")
	if jwksURL == "" {
		log.Fatal("SUPABASE_JWKS_URL is required")
	}
	jwksClient := auth.NewJWKSClient(jwksURL)

	// Repositories
	projectRepo := repository.NewProjectRepository(db)
	nodeRepo := repository.NewNodeRepository(db)
	edgeRepo := repository.NewEdgeRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)

	// Services
	projectService := service.NewProjectService(projectRepo, nodeRepo, edgeRepo)
	nodeService := service.NewNodeService(nodeRepo, edgeRepo)
	edgeService := service.NewEdgeService(edgeRepo)
	settingsService := service.NewSettingsService(settingsRepo)

	// Handlers
	meHandler := handler.NewMeHandler(settingsService)
	projectHandler := handler.NewProjectHandler(projectService)
	nodeHandler := handler.NewNodeHandler(nodeService, projectService)
	edgeHandler := handler.NewEdgeHandler(edgeService, projectService)
	settingsHandler := handler.NewSettingsHandler(settingsService)

	// Router setup
	r := gin.Default()

	// CORS middleware (adjust for production)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	v1 := r.Group("/v1")
	{
		// Auth required routes
		authRequired := v1.Group("")
		authRequired.Use(auth.AuthMiddleware(jwksClient))
		{
			// Me
			authRequired.GET("/me", meHandler.GetMe)

			// Projects
			authRequired.POST("/projects", projectHandler.CreateProject)
			authRequired.GET("/projects", projectHandler.ListProjects)
			authRequired.GET("/projects/:projectId", projectHandler.GetProject)
			authRequired.PATCH("/projects/:projectId", projectHandler.UpdateProject)
			authRequired.GET("/projects/:projectId/tree", projectHandler.GetTree)
			authRequired.POST("/projects/:projectId/save", projectHandler.SaveProject)

			// Nodes
			authRequired.POST("/projects/:projectId/nodes", nodeHandler.CreateNode)
			authRequired.PATCH("/projects/:projectId/nodes/:nodeId", nodeHandler.UpdateNode)
			authRequired.DELETE("/projects/:projectId/nodes/:nodeId", nodeHandler.DeleteNode)

			// Edges
			authRequired.PATCH("/projects/:projectId/edges/:edgeId", edgeHandler.UpdateEdge)
			authRequired.POST("/projects/:projectId/reorder", edgeHandler.Reorder)

			// Settings
			authRequired.GET("/settings", settingsHandler.GetSettings)
			authRequired.PATCH("/settings", settingsHandler.UpdateSettings)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
