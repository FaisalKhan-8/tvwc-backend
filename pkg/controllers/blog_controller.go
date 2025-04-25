package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"your_module/pkg/models"
	"your_module/pkg/services"
	"your_module/pkg/utils"
)

type BlogController struct {
	DB          *gorm.DB
	S3          *utils.S3Utility
	BlogService *services.BlogService
}

func NewBlogController(db *gorm.DB, s3 *utils.S3Utility, blogService *services.BlogService) *BlogController {
	return &BlogController{
		DB:          db,
		S3:          s3,
		BlogService: blogService,
	}
}

func (bc *BlogController) CreateBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog.ID = uuid.New().String()

	err := bc.BlogService.CreateBlog(&blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create blog"})
		return
	}

	c.JSON(http.StatusCreated, blog)
}

func (bc *BlogController) DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	err := bc.BlogService.DeleteBlog(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete blog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "blog deleted successfully"})
}

func (bc *BlogController) UpdateBlog(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	var updatedBlog models.Blog
	if err := c.ShouldBindJSON(&updatedBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := bc.BlogService.UpdateBlog(id, &updatedBlog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update blog"})
		return
	}

	c.JSON(http.StatusOK, updatedBlog)
}

func (bc *BlogController) GetAllBlogs(c *gin.Context) {
	blogs, err := bc.BlogService.GetAllBlogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve blogs"})
		return
	}
	c.JSON(http.StatusOK, blogs)
}

func (bc *BlogController) GetBlogBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}
	blog, err := bc.BlogService.GetBlogBySlug(slug)
	if err != nil {
		if err.Error() == "blog not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to retrieve blog: %v", err)})
		}
		return
	}
	c.JSON(http.StatusOK, blog)
}