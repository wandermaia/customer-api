package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Customer representa a entidade de cliente no sistema
// @Description Entidade que representa um cliente no sistema
type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	Name      string    `json:"name" validate:"required,min=3,max=100" example:"João da Silva"`
	Email     string    `json:"email" validate:"required,email" example:"joao@example.com"`
	Phone     string    `json:"phone" validate:"omitempty,min=8,max=15" example:"(11) 98765-4321"`
	Address   string    `json:"address" validate:"omitempty" example:"Av. Paulista, 1000, São Paulo - SP"`
	Active    bool      `json:"active" gorm:"default:true" example:"true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2025-04-23T15:04:05Z"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2025-04-23T15:04:05Z"`
}

// Validate valida os campos do cliente
func (c *Customer) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
