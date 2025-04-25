package models

type Service struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}