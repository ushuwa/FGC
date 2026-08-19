package models

import "time"

type Client struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	FirstName      string    `gorm:"column:first_name" json:"first_name"`
	LastName       string    `gorm:"column:last_name" json:"last_name"`
	ContactNumber  *string   `gorm:"column:contact_number" json:"contact_number,omitempty"`
	Email          *string   `gorm:"column:email" json:"email,omitempty"`
	CurrentAddress *string   `gorm:"column:current_address" json:"current_address,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Client) TableName() string {
	return "clients"
}
