// models.go
package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	ID          uint        `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name" gorm:"not null"`
	Email       string      `json:"email" gorm:"unique;not null"`
	Preferences Preferences `json:"preferences" gorm:"foreignKey:UserProfileID"`
}

type Preferences struct {
	gorm.Model

	ID               uint   `json:"id" gorm:"primaryKey"`
	UserProfileID    uint   `json:"user_profile_id"`
	TravelStyles     string `json:"travel_styles"`
	PreferredPlaces  string `json:"preferred_places"`
	BudgetRange      string `json:"budget_range"`
	NumberOfNights   int    `json:"number_of_nights"`
	NumberOfRooms    int    `json:"number_of_rooms"`
	NumberOfAdults   int    `json:"number_of_adults"`
	NumberOfChildren int    `json:"number_of_children"`
	FlightType       string `json:"flight_type"`
	DestinationCity  string `json:"destination_city"`
	TravelDate       string `json:"travel_date"` // Use time.Time for actual date handling
}

type Itinerary struct {
	gorm.Model
	ID            uint            `gorm:"primaryKey" json:"id"`
	UserID        uint            `json:"user_id"` // Assuming each itinerary belongs to a user
	TripName      string          `json:"trip_name"`
	StartDate     string          `json:"start_date"` // Use string for JSON compatibility; consider time.Time if you prefer
	EndDate       string          `json:"end_date"`
	Destination   string          `json:"destination"`
	TravelMode    string          `json:"travel_mode"`
	Accommodation string          `json:"accommodation"`
	TotalCost     float64         `json:"total_cost"`
	Activities    string          `json:"activities"`
	Preferences   json.RawMessage `json:"preferences"` // Use json.RawMessage for JSONB support
	Notes         string          `json:"notes"`
	Status        string          `json:"status"`
	Demographics  json.RawMessage `json:"demographics"` // Use json.RawMessage for JSONB support
}
type BookingHistory struct {
	gorm.Model

	ID            uint    `json:"id" gorm:"primaryKey"`
	UserProfileID uint    `json:"user_profile_id"`
	ItineraryID   uint    `json:"itinerary_id"`
	BookingDate   string  `json:"booking_date"` // Use time.Time for actual date handling
	Amount        float64 `json:"amount"`
}
