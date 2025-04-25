package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service struct represents the data model for a service.
type Service struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Image       string             `bson:"image" json:"image"`
	Name        string             `bson:"name" json:"name"`
	Location    string             `bson:"location" json:"location"`
	Description string             `bson:"description" json:"description"`
}