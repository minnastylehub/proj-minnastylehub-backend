package models

// Item represents an item stored in the database
type Item struct {
	ID         string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title      string   `json:"title"`
	Text       string   `json:"text"`
	Brand      string   `json:"brand"`
	Images     []string `json:"images"`
	ButtonLink string   `json:"buttonLink"`
}

type Feedback struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}
