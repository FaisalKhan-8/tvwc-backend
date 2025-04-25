package routes

import (
	"github.com/gin-gonic/gin"
	"your_project_path/pkg/controllers"
	"your_project_path/pkg/middlewares"
)

func ServiceRoutes(router *gin.Engine, serviceController *controllers.ServiceController) {
	serviceGroup := router.Group("/service")
	{
		serviceGroup.POST("", middlewares.AuthMiddleware, serviceController.CreateService)
		serviceGroup.GET("/:id", serviceController.GetService)
		serviceGroup.GET("", serviceController.GetAllServices)
		serviceGroup.DELETE("/:id", middlewares.AuthMiddleware, serviceController.DeleteService)
		serviceGroup.PUT("/:id", middlewares.AuthMiddleware, serviceController.UpdateService)
	}
}