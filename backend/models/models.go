package models

// Pour recevoir le JSON du register
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // ← permet la réception
}

type Post struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Created  string `json:"created"`
}

type ProfileResponse struct {
	User      User     `json:"user"`
	Posts     []Post   `json:"posts"`
	Followers []User   `json:"followers"`
	Following []User   `json:"following"`
	Groups    []string `json:"groups"`
}
