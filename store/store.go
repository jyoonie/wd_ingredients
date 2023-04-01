package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

var ErrNotFound = fmt.Errorf("there's no shit about shit")

type Store interface { //keeping a strict separation between the layers of your service is the biggest benefit of having store interface.
	//So like, your methods of your service shouldn't know anything about your database,
	//they shouldn't rely on a database implementation. Nothing in your service should be dependent on your implementation details.
	//The other benefit of having store interface is certainly that you can switch between databases that fulfill your store interface.
	GetIngredient(ctx context.Context, id uuid.UUID) (*Ingredient, error)
	SearchIngredients(ctx context.Context, i SearchIngredient) ([]Ingredient, error)
	CreateIngredient(ctx context.Context, i Ingredient) (*Ingredient, error)
	UpdateIngredient(ctx context.Context, i Ingredient) (*Ingredient, error)
	DeleteIngredient(ctx context.Context, id uuid.UUID) error

	ListFridgeIngredients(ctx context.Context, i uuid.UUID) ([]FridgeIngredient, error)
	CreateFridgeIngredient(ctx context.Context, f FridgeIngredient) (*FridgeIngredient, error)
	UpdateFridgeIngredient(ctx context.Context, f FridgeIngredient) (*FridgeIngredient, error)
	DeleteFridgeIngredient(ctx context.Context, uid uuid.UUID, iid uuid.UUID) error
}
