package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"your_project/pkg/models"
	"your_project/pkg/services"
	"your_project/pkg/utils"
)

type AboutController struct {
	db             *gorm.DB
	s3Util         *utils.S3Util
	aboutService   *services.AboutService
}

func NewAboutController(db *gorm.DB, s3Util *utils.S3Util, aboutService *services.AboutService) *AboutController {
	return &AboutController{
		db:             db,
		s3Util:         s3Util,
		aboutService: aboutService,
	}
}

func (ac *AboutController) CreateAbout(c *gin.Context) {
	var about models.About
	if err := c.ShouldBindJSON(&about); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	about.ID = uuid.New().String()

	err := ac.aboutService.CreateAbout(&about)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create about"})
		return
	}

	c.JSON(http.StatusCreated, about)
}

func (ac *AboutController) GetAbout(c *gin.Context) {
	about, err := ac.aboutService.GetAbout()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get about"})
		return
	}

	c.JSON(http.StatusOK, about)
}

func (ac *AboutController) UpdateAbout(c *gin.Context) {
	id := c.Param("id")

	var about models.About
	if err := c.ShouldBindJSON(&about); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	about.ID = id
	err := ac.aboutService.UpdateAbout(&about)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update about"})
		return
	}

	c.JSON(http.StatusOK, about)
}

func (ac *AboutController) DeleteAbout(c *gin.Context) {
	id := c.Param("id")

	err := ac.aboutService.DeleteAbout(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete about"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "about deleted"})
}
```
```go
package routes

import (
	"your_project/pkg/controllers"
	"your_project/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAboutRoutes(router *gin.Engine, aboutController *controllers.AboutController) {
	aboutGroup := router.Group("/about")
	{
		aboutGroup.GET("", aboutController.GetAbout)
		aboutGroup.POST("", middlewares.AuthMiddleware(), aboutController.CreateAbout)
		aboutGroup.PUT("/:id", middlewares.AuthMiddleware(), aboutController.UpdateAbout)
		aboutGroup.DELETE("/:id", middlewares.AuthMiddleware(), aboutController.DeleteAbout)
	}
}