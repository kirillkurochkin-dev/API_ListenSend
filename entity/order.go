package entity

import (
	"time"
)

// Order represents an order entity.
type Order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery" gorm:"foreignKey:Email"`
	Payment           Payment   `json:"payment" gorm:"foreignKey:Transaction"`
	Items             []Item    `json:"items" gorm:"foreignKey:TrackNumber"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OffShard          string    `json:"off_shard"`
}

// NewOrder creates a new instance of Order.
func NewOrder() Order {
	return Order{
		// Initialize default values if needed.
	}
}
