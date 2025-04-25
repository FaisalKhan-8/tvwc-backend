package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"pkg/controllers"
	"pkg/routes"
	"pkg/utils"
	"pkg/controllers"
)

func main() {
	fmt.Println("starting server...")

	// Load environment variables.
	port := os.Getenv("PORT")
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatalf("JWT_SECRET not set")
	}
	if os.Getenv("DB_NAME") == "" {
		log.Fatalf("DB_NAME not set")
	}

	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Connect to MongoDB using the utility function
	client, ctx, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database(os.Getenv("DB_NAME"))

	// Dependency Injection
	userController := controllers.NewUserController(db)
	heroController := controllers.NewHeroController(db)
	serviceController := controllers.NewServiceController(db)

	router := gin.Default()

	// Routes Setup
	routes.UserRoutes(router, userController) // Use the imported UserRoutes
	routes.HeroRoutes(router,heroController)
	routes.ServiceRoutes(router, serviceController)

	// Start HTTP Server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}