package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/mokuhyo-driven-test/api/internal/ai"
	"github.com/mokuhyo-driven-test/api/internal/handler"
	"github.com/mokuhyo-driven-test/api/internal/repository"
	postgresRepo "github.com/mokuhyo-driven-test/api/internal/repository/postgres"
	supabaseRepo "github.com/mokuhyo-driven-test/api/internal/repository/supabase"
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

	// DB_TYPEフラグに基づいて実装を選択
	dbType := strings.ToLower(os.Getenv("DB_TYPE"))
	if dbType == "" {
		dbType = "local" // デフォルトはローカルPostgreSQL
	}

	var db repository.DBInterface
	var err error

	switch dbType {
	case "supabase":
		log.Println("Using Supabase database implementation")
		db, err = supabaseRepo.NewSupabaseDB(dbURL)
	case "local", "postgres":
		log.Println("Using local PostgreSQL database implementation")
		db, err = postgresRepo.NewPostgresDB(dbURL)
	default:
		log.Fatalf("Invalid DB_TYPE: %s. Valid values are 'local', 'postgres', or 'supabase'", dbType)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Google OAuth設定
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if googleRedirectURL == "" {
		// フロントエンドのポートを環境変数から取得（デフォルトは3000）
		frontendPort := os.Getenv("FRONTEND_PORT")
		if frontendPort == "" {
			frontendPort = "3000"
		}
		googleRedirectURL = fmt.Sprintf("http://localhost:%s", frontendPort)
	}

	if googleClientID == "" || googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET are required")
	}

	googleOAuthConfig := auth.NewGoogleOAuthConfig(googleClientID, googleClientSecret, googleRedirectURL)
	googleJWKSClient := auth.NewGoogleJWKSClient()

	// Repositories - DB_TYPEに応じて適切な実装を使用
	var projectRepo repository.ProjectRepository
	var nodeRepo repository.NodeRepository
	var edgeRepo repository.EdgeRepository
	var settingsRepo repository.SettingsRepository
	var userRepo repository.UserRepository

	switch dbType {
	case "supabase":
		projectRepo = supabaseRepo.NewProjectRepository(db)
		nodeRepo = supabaseRepo.NewNodeRepository(db)
		edgeRepo = supabaseRepo.NewEdgeRepository(db)
		settingsRepo = supabaseRepo.NewSettingsRepository(db)
		userRepo = supabaseRepo.NewUserRepository(db)
	case "local", "postgres":
		projectRepo = postgresRepo.NewProjectRepository(db)
		nodeRepo = postgresRepo.NewNodeRepository(db)
		edgeRepo = postgresRepo.NewEdgeRepository(db)
		settingsRepo = postgresRepo.NewSettingsRepository(db)
		userRepo = postgresRepo.NewUserRepository(db)
	}

	// Services
	authService := service.NewAuthService(userRepo)
	projectService := service.NewProjectService(projectRepo, nodeRepo, edgeRepo)

	var questionSelector ai.QuestionSelector
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiAPIKey != "" {
		selector, err := ai.NewGeminiQuestionSelector(context.Background(), geminiAPIKey, geminiModel)
		if err != nil {
			log.Printf("Gemini client init failed: %v", err)
		} else {
			questionSelector = selector
			defer selector.Close()
		}
	} else {
		log.Println("GEMINI_API_KEY is not set, using fallback question selection")
	}

	nodeService := service.NewNodeService(nodeRepo, edgeRepo, questionSelector)
	edgeService := service.NewEdgeService(edgeRepo)
	settingsService := service.NewSettingsService(settingsRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, googleOAuthConfig)
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
		// Public auth routes
		v1.POST("/auth/google", authHandler.HandleGoogleAuth)

		// Auth required routes
		authRequired := v1.Group("")
		// Google認証ミドルウェアを使用
		authRequired.Use(func(c *gin.Context) {
			// ミドルウェア内でユーザーIDを取得する関数
			getUserIDByGoogleID := func(googleID string) (uuid.UUID, error) {
				user, err := authService.GetUserByGoogleID(c.Request.Context(), googleID)
				if err != nil || user == nil {
					return uuid.Nil, fmt.Errorf("user not found")
				}
				return user.ID, nil
			}
			// Google認証ミドルウェアを実行
			middleware := auth.GoogleAuthMiddleware(googleJWKSClient, getUserIDByGoogleID)
			middleware(c)
		})
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
