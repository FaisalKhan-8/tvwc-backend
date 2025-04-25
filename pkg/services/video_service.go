package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"your_project/pkg/middlewares"
	"your_project/pkg/models"
	"your_project/pkg/utils"
)

type VideoService struct {
	DB    *gorm.DB
	S3    *utils.S3Client
	Store string
}

func NewVideoService(db *gorm.DB, s3 *utils.S3Client) *VideoService {
	return &VideoService{
		DB:    db,
		S3:    s3,
		Store: os.Getenv("AWS_BUCKET"),
	}
}

func (s *VideoService) CreateVideo(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "User is not admin", http.StatusUnauthorized)
		return
	}

	var video models.Video
	if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error getting the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}
	extension := filepath.Ext(header.Filename)
	uuidValue := uuid.New()
	key := fmt.Sprintf("%s%s", uuidValue, extension)
	uploadUrl, err := s.S3.UploadFile(key, fileBytes, s.Store)
	if err != nil {
		http.Error(w, "Error while uploading the file", http.StatusInternalServerError)
		return
	}
	video.VideoURL = uploadUrl
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()

	result := s.DB.Create(&video)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

func (s *VideoService) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "User is not admin", http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var video models.Video
	if err := s.DB.First(&video, id).Error; err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	key := video.VideoURL[len("https://"+s.Store+".s3.amazonaws.com/"):]
	if err := s.S3.DeleteFile(key, s.Store); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := s.DB.Delete(&video)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Video deleted successfully"})
}

func (s *VideoService) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "User is not admin", http.StatusUnauthorized)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var video models.Video
	if err := s.DB.First(&video, id).Error; err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	var updatedVideo models.Video
	if err := json.NewDecoder(r.Body).Decode(&updatedVideo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("video")
	if err != nil && file != nil {
		http.Error(w, "Error getting the file", http.StatusBadRequest)
		return
	}
	if file != nil {
		defer file.Close()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading the file", http.StatusInternalServerError)
			return
		}
		key := video.VideoURL[len("https://"+s.Store+".s3.amazonaws.com/"):]
		if err := s.S3.DeleteFile(key, s.Store); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		extension := filepath.Ext(header.Filename)
		uuidValue := uuid.New()
		key = fmt.Sprintf("%s%s", uuidValue, extension)
		uploadUrl, err := s.S3.UploadFile(key, fileBytes, s.Store)
		if err != nil {
			http.Error(w, "Error while uploading the file", http.StatusInternalServerError)
			return
		}
		video.VideoURL = uploadUrl
	}

	video.Category = updatedVideo.Category
	video.Title = updatedVideo.Title
	video.Content = updatedVideo.Content
	video.UpdatedAt = time.Now()

	result := s.DB.Save(&video)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(video)
}

func (s *VideoService) GetVideo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var video models.Video
	if err := s.DB.First(&video, id).Error; err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(video)
}

func (s *VideoService) GetPublicVideos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	result := s.DB.Find(&videos)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videos)
}