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
	"pkg/services"
	
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

	// Connect to MongoDB using the utility function from db.go
	client, ctx, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database(os.Getenv("DB_NAME")) // Use the global db variable
	// We don't use database migrations for this project because we're using MongoDB.
	// MongoDB is a NoSQL database that doesn't require schema migrations in the same way as relational databases.
	// Schema changes can be handled dynamically within the application.

	// Dependency Injection
	userController := controllers.NewUserController(db)
	heroController := controllers.NewHeroController(db)
	serviceController := controllers.NewServiceController(db)
	videoService := services.NewVideoService(db, utils.NewS3Client())
	videoController := controllers.NewVideoController(videoService)
	aboutService := services.NewAboutService(db, utils.NewS3Client())
	aboutController := controllers.NewAboutController(aboutService)
	blogService := services.NewBlogService(db, utils.NewS3Client())
	blogController := controllers.NewBlogController(blogService)

	router := gin.Default()

	// Routes Setup
	routes.VideoRoutes(router, videoController)
	routes.BlogRoutes(router, blogController)	

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
	routes.AboutRoutes(router, aboutController)

}