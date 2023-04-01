package store

import (
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	IngredientUUID uuid.UUID
	IngredientName string
	Category       string
	DaysUntilExp   int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// type Query struct {
// 	IngredientName *string //I have to check if they're null. If they're not null, I'm going to add each query to base sql query
// 	Category       *string
// 	DaysUntilExp   *int
// }

type SearchIngredient struct {
	IngredientName *string
	Category       *string
}

type FridgeIngredient struct {
	UserUUID       uuid.UUID
	IngredientUUID uuid.UUID
	Amount         int
	Unit           string
	PurchasedDate  time.Time
	ExpirationDate time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
