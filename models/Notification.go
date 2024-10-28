package models

type Notification struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
	Status  string `json:"status"` // e.g., "sent", "pending"
}
