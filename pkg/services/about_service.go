package services

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/yourusername/yourproject/pkg/models"
	"github.com/yourusername/yourproject/pkg/utils"
)

// AboutService struct to hold dependencies
type AboutService struct {
	DB    *gorm.DB
	S3Util *utils.S3Util
}

// NewAboutService creates a new AboutService instance
func NewAboutService(db *gorm.DB, s3Util *utils.S3Util) *AboutService {
	return &AboutService{DB: db, S3Util: s3Util}
}

// CreateAbout creates a new about entry
func (s *AboutService) CreateAbout(about *models.About) error {
	result := s.DB.Create(about)
	if result.Error != nil {
		return fmt.Errorf("failed to create about: %w", result.Error)
	}
	return nil
}

// GetAbout returns the about entry
func (s *AboutService) GetAbout() (*models.About, error) {
	var abouts []models.About
	result := s.DB.Find(&abouts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get about: %w", result.Error)
	}
	if len(abouts) == 0 {
		return nil, errors.New("no about found")
	}
	return &abouts[0], nil
}

// UpdateAbout updates an existing about entry
func (s *AboutService) UpdateAbout(about *models.About) error {
	result := s.DB.Save(about)
	if result.Error != nil {
		return fmt.Errorf("failed to update about: %w", result.Error)
	}
	return nil
}

// DeleteAbout deletes an about entry by ID
func (s *AboutService) DeleteAbout(id string) error {
	result := s.DB.Delete(&models.About{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete about: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("no about found with given id")
	}
	return nil
}