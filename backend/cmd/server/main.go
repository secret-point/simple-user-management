package main

import (
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"context"
	"time"
	
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"arritech-challenge/pkg/db"        
	"arritech-challenge/pkg/handlers"  
	"arritech-challenge/pkg/middlewares"
)

func main() {
	dbClient := db.ConnectDB()
	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Error("Error on disconnection with MongoDB: %v", err)
		}
	}()

	// Initialize Gin router
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.Use(middlewares.RequestLogger())

	// Setting up routes with the new approach
	userHandler := handlers.NewUserHandler(dbClient)
	router.GET("/users", userHandler.GetUsers)
	router.POST("/user", userHandler.CreateUser)
	router.GET("/user/:id", userHandler.GetUser)
	router.PUT("/user/:id", userHandler.UpdateUser)
	router.DELETE("/user/:id", userHandler.DeleteUser)

	// Run server with graceful shutdown
	runServer(router)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func runServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown:", err)
	}

	log.Info("Server exiting")
}
