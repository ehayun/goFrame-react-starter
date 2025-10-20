package model

// AppResource matches the existing app_resources table in the database
type AppResource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
