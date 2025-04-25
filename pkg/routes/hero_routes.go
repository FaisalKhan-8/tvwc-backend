package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"your-module-path/pkg/controllers"
	"your-module-path/pkg/middlewares"
)

func HeroRoutes(router *gin.Engine, heroController controllers.HeroController) {
	heroGroup := router.Group("/hero")
	{
		heroGroup.GET("/", func(c *gin.Context) {
			hero, err := heroController.GetHeroSection(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, hero)
		})

		adminGroup := heroGroup.Group("/", middlewares.AdminAuthMiddleware())
		{
			adminGroup.POST("/", func(c *gin.Context) {
				err := heroController.CreateHeroSection(c)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"message": "Hero section created successfully"})
			})

			adminGroup.DELETE("/:id", func(c *gin.Context) {
				err := heroController.DeleteHeroSection(c)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Hero section deleted successfully"})
			})

			adminGroup.PUT("/:id", func(c *gin.Context) {
				err := heroController.UpdateHeroSection(c)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Hero section updated successfully"})
			})
		}
	}
}