package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"yourapp/pkg/middlewares"
	"yourapp/pkg/models"
	"yourapp/pkg/services"
	"yourapp/pkg/utils"
)

type VideoController struct {
	VideoService *services.VideoService
	S3Utils      *utils.S3Utils
	DBUtils      *utils.DBUtils
}

func NewVideoController(videoService *services.VideoService, s3Utils *utils.S3Utils, dbUtils *utils.DBUtils) *VideoController {
	return &VideoController{
		VideoService: videoService,
		S3Utils:      s3Utils,
		DBUtils:      dbUtils,
	}
}

func (vc *VideoController) CreateVideo(w http.ResponseWriter, r *http.Request) {
	middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var video models.Video
		err := json.NewDecoder(r.Body).Decode(&video)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		video.CreatedAt = time.Now()
		video.UpdatedAt = time.Now()

		err = vc.VideoService.CreateVideo(&video)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(video)
	})).ServeHTTP(w, r)
}

func (vc *VideoController) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid video ID", http.StatusBadRequest)
			return
		}

		err = vc.VideoService.DeleteVideo(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(w, r)
}

func (vc *VideoController) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid video ID", http.StatusBadRequest)
			return
		}

		var video models.Video
		err = json.NewDecoder(r.Body).Decode(&video)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		video.ID = id
		video.UpdatedAt = time.Now()

		err = vc.VideoService.UpdateVideo(&video)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(video)
	})).ServeHTTP(w, r)
}

func (vc *VideoController) GetVideo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	video, err := vc.VideoService.GetVideo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(video)
}

func (vc *VideoController) GetAllPublicVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := vc.VideoService.GetAllPublicVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(videos)
}

func (vc *VideoController) RegisterVideoRoutes(router *mux.Router) {
	router.HandleFunc("/videos", vc.CreateVideo).Methods("POST")
	router.HandleFunc("/videos/{id}", vc.DeleteVideo).Methods("DELETE")
	router.HandleFunc("/videos/{id}", vc.UpdateVideo).Methods("PUT")
	router.HandleFunc("/videos/{id}", vc.GetVideo).Methods("GET")
	router.HandleFunc("/videos", vc.GetAllPublicVideos).Methods("GET")
}
```
```go
package services

import (
	"fmt"
	"yourapp/pkg/models"
	"yourapp/pkg/utils"
)

type VideoService struct {
	DBUtils *utils.DBUtils
}

func NewVideoService(dbUtils *utils.DBUtils) *VideoService {
	return &VideoService{
		DBUtils: dbUtils,
	}
}

func (vs *VideoService) CreateVideo(video *models.Video) error {
	result := vs.DBUtils.DB.Create(video)
	if result.Error != nil {
		return fmt.Errorf("failed to create video: %w", result.Error)
	}
	return nil
}

func (vs *VideoService) DeleteVideo(id int) error {
	result := vs.DBUtils.DB.Delete(&models.Video{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete video: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("video not found with ID %d", id)
	}
	return nil
}

func (vs *VideoService) UpdateVideo(video *models.Video) error {
	result := vs.DBUtils.DB.Save(video)
	if result.Error != nil {
		return fmt.Errorf("failed to update video: %w", result.Error)
	}
	return nil
}

func (vs *VideoService) GetVideo(id int) (*models.Video, error) {
	var video models.Video
	result := vs.DBUtils.DB.First(&video, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get video: %w", result.Error)
	}
	return &video, nil
}

func (vs *VideoService) GetAllPublicVideos() ([]models.Video, error) {
	var videos []models.Video
	result := vs.DBUtils.DB.Find(&videos)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get videos: %w", result.Error)
	}
	return videos, nil
}