package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"your_module_name/pkg/models"
)

// HeroController struct to hold dependencies
type HeroController struct {
	collection *mongo.Collection
	ctx        context.Context
}

// NewHeroController creates a new HeroController
func NewHeroController(collection *mongo.Collection, ctx context.Context) *HeroController {
	return &HeroController{
		collection: collection,
		ctx:        ctx,
	}
}

// CreateHero creates a new hero section
func (hc *HeroController) CreateHero(c *gin.Context) {
	var hero models.Hero
	if err := c.ShouldBindJSON(&hero); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hero.ID = primitive.NewObjectID()
	hero.CreatedAt = time.Now()
	hero.UpdatedAt = time.Now()

	_, err := hc.collection.InsertOne(hc.ctx, hero)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hero section"})
		return
	}

	c.JSON(http.StatusCreated, hero)
}

// GetHero gets the hero section
func (hc *HeroController) GetHero(c *gin.Context) {
	var hero models.Hero
	err := hc.collection.FindOne(hc.ctx, bson.M{}).Decode(&hero)
	if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Hero section not found"})
            return
        }
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get hero section"})
		return
	}

	c.JSON(http.StatusOK, hero)
}

// UpdateHero updates the hero section
func (hc *HeroController) UpdateHero(c *gin.Context) {
	var hero models.Hero
	if err := c.ShouldBindJSON(&hero); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hero.UpdatedAt = time.Now()

	_, err := hc.collection.UpdateOne(hc.ctx, bson.M{}, bson.M{"$set": hero})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hero section"})
		return
	}

	c.JSON(http.StatusOK, hero)
}

// DeleteHero deletes the hero section
func (hc *HeroController) DeleteHero(c *gin.Context) {
	_, err := hc.collection.DeleteOne(hc.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hero section"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hero section deleted successfully"})
}