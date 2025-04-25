package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/pkg/controllers"
	"github.com/yourusername/yourproject/pkg/middlewares"
)

func VideoRoutes(router *gin.Engine, videoController *controllers.VideoController) {
	videoGroup := router.Group("/videos")
	{
		videoGroup.GET("", videoController.GetAllPublicVideos)
	}

	adminVideoGroup := router.Group("/admin/videos")
	adminVideoGroup.Use(middlewares.AuthMiddleware())
	{
		adminVideoGroup.POST("", videoController.CreateVideo)
		adminVideoGroup.GET("/:id", videoController.GetVideo)
		adminVideoGroup.PUT("/:id", videoController.UpdateVideo)
		adminVideoGroup.DELETE("/:id", videoController.DeleteVideo)
	}
}