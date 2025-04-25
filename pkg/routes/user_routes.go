package routes

import (
	"your_project_name/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userController *controllers.UserController) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/signup", userController.Signup)
		userGroup.POST("/login", userController.Login)
		//userGroup.GET("/:id", userController.GetUser)
		//userGroup.PUT("/:id", userController.UpdateUser)
		//userGroup.DELETE("/:id", userController.DeleteUser)
	}
}