package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/pkg/controllers"
	"github.com/yourusername/yourproject/pkg/middlewares"
)

func BlogRoutes(router *gin.Engine, blogController *controllers.BlogController) {
	blogGroup := router.Group("/api/blogs")
	{
		blogGroup.GET("", blogController.GetAllBlogs)
		blogGroup.GET("/:slug", blogController.GetBlogBySlug)

		adminBlogGroup := blogGroup.Group("", middlewares.AuthMiddleware())
		{
			adminBlogGroup.POST("", blogController.CreateBlog)
			adminBlogGroup.PUT("/:id", blogController.UpdateBlog)
			adminBlogGroup.DELETE("/:id", blogController.DeleteBlog)
		}
	}
}