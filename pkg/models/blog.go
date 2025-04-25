package models

import "time"

type Blog struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Content         string    `json:"content"`
	ImageURL        string    `json:"image_url"`
	Author          string    `json:"author"`
	AuthorImageURL  string    `json:"author_image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}