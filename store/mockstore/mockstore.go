package mockstore

import (
	"context"
	"time"
	"wd_ingredients/store"

	"github.com/google/uuid"
)

var _ store.Store = (*Mockstore)(nil)

type Mockstore struct {
	GetIngredientOverride     func(ctx context.Context, id uuid.UUID) (*store.Ingredient, error)
	SearchIngredientsOverride func(ctx context.Context, i store.SearchIngredient) ([]store.Ingredient, error)
	CreateIngredientOverride  func(ctx context.Context, i store.Ingredient) (*store.Ingredient, error)
	UpdateIngredientOverride  func(ctx context.Context, i store.Ingredient) (*store.Ingredient, error)
	DeleteIngredientOverride  func(ctx context.Context, id uuid.UUID) error

	ListFridgeIngredientsOverride  func(ctx context.Context, id uuid.UUID) ([]store.FridgeIngredient, error)
	CreateFridgeIngredientOverride func(ctx context.Context, f store.FridgeIngredient) (*store.FridgeIngredient, error)
	UpdateFridgeIngredientOverride func(ctx context.Context, f store.FridgeIngredient) (*store.FridgeIngredient, error)
	DeleteFridgeIngredientOverride func(ctx context.Context, uid uuid.UUID, iid uuid.UUID) error
}

func (m *Mockstore) GetIngredient(ctx context.Context, id uuid.UUID) (*store.Ingredient, error) {
	if m.GetIngredientOverride != nil {
		return m.GetIngredientOverride(ctx, id)
	}
	return &store.Ingredient{
		IngredientUUID: id,
		IngredientName: "onion",
		Category:       "vegetables",
		DaysUntilExp:   7,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (m *Mockstore) SearchIngredients(ctx context.Context, i store.SearchIngredient) ([]store.Ingredient, error) {
	if m.SearchIngredientsOverride != nil {
		return m.SearchIngredientsOverride(ctx, i)
	}

	return []store.Ingredient{
		{
			IngredientUUID: uuid.MustParse("080b5f09-527b-4581-bb56-19adbfe50ebf"),
			IngredientName: "tuna",
			Category:       "fish",
			DaysUntilExp:   3,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		// {
		// 	IngredientUUID: uuid.MustParse("080b5f09-527b-4581-bb56-19adbfe50ebf"),
		// 	IngredientName: "tuna",
		// 	Category:       "tuna sushi",
		// 	DaysUntilExp:   3,
		// 	CreatedAt:      time.Now(),
		// 	UpdatedAt:      time.Now(),
		// },
	}, nil
}

func (m *Mockstore) CreateIngredient(ctx context.Context, i store.Ingredient) (*store.Ingredient, error) {
	if m.CreateIngredientOverride != nil {
		return m.CreateIngredientOverride(ctx, i)
	}

	return &store.Ingredient{
		IngredientUUID: uuid.New(),
		IngredientName: i.IngredientName,
		Category:       i.Category,
		DaysUntilExp:   i.DaysUntilExp,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (m *Mockstore) UpdateIngredient(ctx context.Context, i store.Ingredient) (*store.Ingredient, error) {
	if m.UpdateIngredientOverride != nil {
		return m.UpdateIngredientOverride(ctx, i)
	}

	i.UpdatedAt = time.Now()

	return &i, nil
}

func (m *Mockstore) DeleteIngredient(ctx context.Context, id uuid.UUID) error {
	if m.DeleteIngredientOverride != nil {
		return m.DeleteIngredientOverride(ctx, id)
	}

	return nil
}

func (m *Mockstore) ListFridgeIngredients(ctx context.Context, id uuid.UUID) ([]store.FridgeIngredient, error) {
	if m.ListFridgeIngredientsOverride != nil {
		return m.ListFridgeIngredientsOverride(ctx, id)
	}

	return []store.FridgeIngredient{
		{
			UserUUID:       id,
			IngredientUUID: uuid.MustParse("ffff7c73-52b0-4e3d-bf3f-0c26785ef972"),
			Amount:         3,
			Unit:           "kg",
			PurchasedDate:  time.Date(2023, time.March, 24, 15, 0, 0, 0, time.Now().Location()),
			ExpirationDate: time.Date(2023, time.March, 24, 15, 0, 0, 0, time.Now().Location()),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			UserUUID:       id,
			IngredientUUID: uuid.MustParse("2c98fff4-7ccc-4536-8259-67a88380e99c"),
			Amount:         2,
			Unit:           "L",
			PurchasedDate:  time.Date(2023, time.March, 24, 15, 0, 0, 0, time.Now().Location()),
			ExpirationDate: time.Date(2023, time.March, 31, 15, 0, 0, 0, time.Now().Location()),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}, nil
}

func (m *Mockstore) CreateFridgeIngredient(ctx context.Context, f store.FridgeIngredient) (*store.FridgeIngredient, error) {
	if m.CreateFridgeIngredientOverride != nil {
		return m.CreateFridgeIngredientOverride(ctx, f)
	}

	f.CreatedAt = time.Now() //time.Date(2023, time.March, 24, 5, 5, 5, 5, time.Now().UTC().Location()) 이런식으로 아무렇게나 써도 사실 핸들러 테스트에서 테스트하는건 api 모델이므로 이것도 통과함..
	f.UpdatedAt = time.Now()

	// return &store.FridgeIngredient{
	// 	UserUUID:       f.UserUUID,
	// 	IngredientUUID: f.IngredientUUID,
	// 	Amount:         f.Amount,
	// 	Unit:           f.Unit,
	// 	PurchasedDate:  f.PurchasedDate,
	// 	ExpirationDate: f.ExpirationDate,
	// 	CreatedAt:      time.Now(),
	// 	UpdatedAt:      time.Now(),
	// }
	return &f, nil
}

func (m *Mockstore) UpdateFridgeIngredient(ctx context.Context, f store.FridgeIngredient) (*store.FridgeIngredient, error) {
	if m.UpdateFridgeIngredientOverride != nil {
		return m.UpdateFridgeIngredientOverride(ctx, f)
	}

	f.UpdatedAt = time.Now()

	return &f, nil
}

func (m *Mockstore) DeleteFridgeIngredient(ctx context.Context, uid uuid.UUID, fid uuid.UUID) error {
	if m.DeleteFridgeIngredientOverride != nil {
		return m.DeleteFridgeIngredientOverride(ctx, uid, fid)
	}

	return nil
}
