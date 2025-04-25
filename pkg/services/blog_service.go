package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"awesomeProject/pkg/models"
	"awesomeProject/pkg/utils"
)

type BlogService struct {
	DB      *gorm.DB
	S3Utils utils.S3Utils
}

func NewBlogService(db *gorm.DB, s3Utils utils.S3Utils) *BlogService {
	return &BlogService{
		DB:      db,
		S3Utils: s3Utils,
	}
}

func (s *BlogService) CreateBlog(blog *models.Blog) error {
	if blog.Slug == "" {
		blog.Slug = generateSlug(blog.Title)
	}
	if err := s.DB.Create(blog).Error; err != nil {
		return fmt.Errorf("failed to create blog: %w", err)
	}
	return nil
}

func (s *BlogService) UpdateBlog(id string, blog *models.Blog) error {
	var existingBlog models.Blog
	if err := s.DB.First(&existingBlog, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("blog not found")
		}
		return fmt.Errorf("failed to find blog: %w", err)
	}
	if blog.Slug == "" {
		blog.Slug = generateSlug(blog.Title)
	}

	blog.ID = existingBlog.ID

	if err := s.DB.Save(blog).Error; err != nil {
		return fmt.Errorf("failed to update blog: %w", err)
	}
	return nil
}

func (s *BlogService) DeleteBlog(id string) error {
	var blog models.Blog
	if err := s.DB.First(&blog, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("blog not found")
		}
		return fmt.Errorf("failed to find blog: %w", err)
	}

	if err := s.DB.Delete(&blog).Error; err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	return nil
}

func (s *BlogService) GetBlog(id string) (*models.Blog, error) {
	var blog models.Blog
	if err := s.DB.First(&blog, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("blog not found")
		}
		return nil, fmt.Errorf("failed to find blog: %w", err)
	}
	return &blog, nil
}

func (s *BlogService) GetAllBlogs() ([]models.Blog, error) {
	var blogs []models.Blog
	if err := s.DB.Find(&blogs).Error; err != nil {
		return nil, fmt.Errorf("failed to get blogs: %w", err)
	}
	return blogs, nil
}

func (s *BlogService) GetBlogBySlug(slug string) (*models.Blog, error) {
	var blog models.Blog
	if err := s.DB.Where("slug = ?", slug).First(&blog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("blog with slug %s not found", slug)
		}
		return nil, fmt.Errorf("failed to get blog by slug: %w", err)
	}
	return &blog, nil
}

func generateSlug(title string) string {
	return fmt.Sprintf("%s-%s", utils.Slugify(title), uuid.New().String()[:8])
}