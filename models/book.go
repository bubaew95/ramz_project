package models

type Book struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Year    int    `json:"year"`
	Image   string `json:"image"`
	Visible int    `json:"visible"`
	UserID  *int   `json:"user_id,omitempty"`
}
