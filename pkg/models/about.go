package models

type About struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Subtitle         string `json:"subtitle"`
	Description      string `json:"description"`
	ImageURL         string `json:"image_url"`
	YearsExperience  string `json:"years_experience"`
	ProjectChallenge string `json:"project_challenge"`
	PositiveReviews  string `json:"positive_reviews"`
	TrustedStudents  string `json:"trusted_students"`
}