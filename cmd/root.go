package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/raihankhan/jwt-auth-system/config"
	"github.com/raihankhan/jwt-auth-system/internal/database"
	"github.com/raihankhan/jwt-auth-system/internal/handler"
	"github.com/raihankhan/jwt-auth-system/internal/middleware"
	"github.com/raihankhan/jwt-auth-system/internal/user"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jwt-auth-system",
	Short: "JWT Authentication and Authorization System",
	Long: `This application demonstrates a JWT-based authentication and authorization system
using Go, Gin, Cobra, and other technologies.`,
	Run: startServer(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/config.yaml", "config file (default is config/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggleDebug", "d", false, "Enable debug logging")

	// Initialize migrate commands
	SetupMigrateCommands(rootCmd) // Call SetupMigrateCommands from cmd/migrate.go
}

func startServer() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
			os.Exit(1)
		}

		// Database Connection
		db, err := database.ConnectDB(cfg)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
			os.Exit(1)
		}
		defer func() {
			sqlDB, _ := db.DB()
			sqlDB.Close()
			log.Println("Database connection closed.")
		}()

		// Database Migration
		// ref: https://gorm.io/docs/migration.html
		log.Println("Running database migrations...")
		err = db.AutoMigrate(&user.User{})
		if err != nil {
			log.Fatalf("Failed to run database migrations: %v", err)
			os.Exit(1)
		}
		log.Println("Database migrations completed successfully.")

		// Initialize Gin
		router := gin.Default()

		// Middleware to make config available in context
		router.Use(func(c *gin.Context) { // Add middleware function
			c.Set("config", cfg) // Set config in Gin context
			c.Next()             // Proceed to the next handler
		})

		// Initialize User Handler
		userHandler := handler.NewUserHandler(db)

		// Define API routes
		api := router.Group("/api")
		{
			api.POST("/register", userHandler.RegisterUser)
			api.POST("/login", userHandler.LoginUser) // Login endpoint

			// Protected endpoint - requires authentication
			api.GET("/protected", middleware.AuthMiddleware(), userHandler.ProtectedEndpoint) // Apply AuthMiddleware
		}

		// Start the Gin server
		serverAddress := fmt.Sprintf(":%s", cfg.App.Port)
		log.Printf("Starting server on %s...", serverAddress)
		if err := router.Run(serverAddress); err != nil {
			log.Fatalf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}
}
