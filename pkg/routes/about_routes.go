package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-alura-challenge-backend/pkg/controllers"
	"github.com/jtonynet/api-alura-challenge-backend/pkg/middlewares"
)

func AboutRoutes(router *gin.Engine, aboutController *controllers.AboutController) {
	about := router.Group("/api/about")
	{
		about.GET("", aboutController.GetAbout)

		aboutAuth := about.Group("")
		aboutAuth.Use(middlewares.AuthMiddleware())
		{
			aboutAuth.POST("", aboutController.CreateAbout)
			aboutAuth.PUT("/:id", aboutController.UpdateAbout)
			aboutAuth.DELETE("/:id", aboutController.DeleteAbout)
		}
	}
}