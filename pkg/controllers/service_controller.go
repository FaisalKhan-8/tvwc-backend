package controllers

import (
	"context"
	"net/http"

	"your_project_name/pkg/models"
	"your_project_name/pkg/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServiceController handles service-related operations.
type ServiceController struct {
	ServiceService *services.ServiceService
}

// NewServiceController creates a new ServiceController.
func NewServiceController(serviceService *services.ServiceService) *ServiceController {
	return &ServiceController{
		ServiceService: serviceService,
	}
}

// CreateService creates a new service.
// @Summary Create a new service
// @Description Create a new service with the given data.
// @Tags services
// @Accept json
// @Produce json
// @Param service body models.Service true "Service data"
// @Success 201 {object} models.Service
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service [post]
func (sc *ServiceController) CreateService(c *gin.Context) {
	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdService, err := sc.ServiceService.CreateService(context.TODO(), service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdService)
}

// GetService retrieves a service by ID.
// @Summary Get a service by ID
// @Description Get a service by its ID.
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} models.Service
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service/{id} [get]
func (sc *ServiceController) GetService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	service, err := sc.ServiceService.GetService(context.TODO(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, service)
}

// GetAllServices retrieves all services.
// @Summary Get all services
// @Description Get all services.
// @Tags services
// @Produce json
// @Success 200 {array} models.Service
// @Failure 500 {object} map[string]string
// @Router /service [get]
func (sc *ServiceController) GetAllServices(c *gin.Context) {
	services, err := sc.ServiceService.GetAllServices(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// UpdateService updates an existing service.
// @Summary Update a service
// @Description Update an existing service with the given data.
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body models.Service true "Updated service data"
// @Success 200 {object} models.Service
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service/{id} [put]
func (sc *ServiceController) UpdateService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.ID = id

	updatedService, err := sc.ServiceService.UpdateService(context.TODO(), service)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedService)
}

// DeleteService deletes a service by ID.
// @Summary Delete a service
// @Description Delete a service by its ID.
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service/{id} [delete]
func (sc *ServiceController) DeleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = sc.ServiceService.DeleteService(context.TODO(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}