package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	NAME     string `json:"name"`
	EMAIL    string `json:"email"`
	PASSWORD string `json:"password"`
}

type Preferences struct {
	gorm.Model
	UserID        uint            `json:"user_id"`
	LOCATIONS     json.RawMessage `json:"locations"`
	BUDGET        string          `json:"budget"`
	BEDROOMS      uint            `json:"bedrooms"`
	PROPERTY_SIZE float64         `json:"property_size"`
	AMENITIES     json.RawMessage `json:"amenities"`

	User User `gorm:"foreignKey:UserID"`
}

type Booking struct {
	gorm.Model
	PropertyID   uint   `json:"property_id"`
	BookingDate  string `json:"booking_date"`
	BookingTime  string `json:"booking_time"`
	CheckoutDate string `json:"checkout_date"`
	CheckoutTime string `json:"checkout_time"`
	UserID       uint   `json:"user_id"`

	User     User     `gorm:"foreignKey:UserID"`
	Property Property `gorm:"foreignKey:PropertyID"`
}

type Property struct {
	gorm.Model
	Title         string          `gorm:"size:500" json:"title"`
	Description   string          `gorm:"type:text" json:"description"`
	PropertyType  string          `gorm:"type:varchar(100);not null" json:"property_type"`
	Address       string          `gorm:"size:500" json:"address"`
	City          string          `gorm:"size:200" json:"city"`
	Price         float64         `json:"price"`
	Currency      string          `gorm:"size:10;default:USD" json:"currency"`
	PricePeriod   string          `gorm:"size:50" json:"price_period"`
	Bedrooms      uint            `json:"bedrooms"`
	Bathrooms     uint            `json:"bathrooms"`
	AreaSqft      float64         `json:"area_sqft"`
	Amenities     json.RawMessage `json:"amenities"` // JSON string
	SourceWebsite string          `gorm:"size:200;not null" json:"source_website"`
	SourceURL     string          `gorm:"size:1000" json:"source_url"`
	ExternalID    string          `gorm:"size:200" json:"external_id"`
	ImageURLs     string          `gorm:"type:text" json:"image_urls"` // JSON string
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	LastScrapedAt time.Time       `json:"last_scraped_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}
