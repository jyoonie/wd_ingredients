package service

import (
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	IngredientUUID uuid.UUID `json:"ingredient_uuid,omitempty"`
	IngredientName string    `json:"ingredient_name,omitempty"`
	Category       string    `json:"category,omitempty"`
	DaysUntilExp   int       `json:"days_until_exp,omitempty"`
	//created_at           time.Time
	//updated_at           time.Time
}

type SearchIngredients struct {
	IngredientName string `json:"ingredient_name,omitempty"`
	Category       string `json:"category,omitempty"`
}

type FridgeIngredient struct {
	UserUUID       uuid.UUID `json:"user_uuid,omitempty"`
	IngredientUUID uuid.UUID `json:"ingredient_uuid,omitempty"`
	Amount         int       `json:"amount,omitempty"`
	Unit           string    `json:"unit,omitempty"`
	PurchasedDate  time.Time `json:"purchased_date,omitempty"`
	ExpirationDate time.Time `json:"expiration_date,omitempty"`
	// created_at       time.Time
	// updated_at       time.Time
}
