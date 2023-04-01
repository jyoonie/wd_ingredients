package service

import (
	"wd_ingredients/store"
)

func apiIngr2DBIngr(i Ingredient) store.Ingredient {
	return store.Ingredient{
		IngredientUUID: i.IngredientUUID,
		IngredientName: i.IngredientName,
		Category:       i.Category,
		DaysUntilExp:   i.DaysUntilExp,
	}
}

func dbIngr2ApiIngr(i *store.Ingredient) Ingredient {
	return Ingredient{
		IngredientUUID: i.IngredientUUID,
		IngredientName: i.IngredientName,
		Category:       i.Category,
		DaysUntilExp:   i.DaysUntilExp,
	}
}

func apiSearchIngr2DBSearchIngr(i SearchIngredients) store.SearchIngredient {
	var out store.SearchIngredient

	if i.IngredientName != "" {
		out.IngredientName = &i.IngredientName
	}
	if i.Category != "" {
		out.Category = &i.Category
	}
	return out
}

func apiFIngr2DBFIngr(f FridgeIngredient) store.FridgeIngredient { //api model에 있는 필드만 신경써.
	return store.FridgeIngredient{
		UserUUID:       f.UserUUID,
		IngredientUUID: f.IngredientUUID,
		Amount:         f.Amount,
		Unit:           f.Unit,
		PurchasedDate:  f.PurchasedDate,
		ExpirationDate: f.ExpirationDate,
	}
}

func dbFIngr2ApiFIngr(f *store.FridgeIngredient) FridgeIngredient {
	return FridgeIngredient{
		UserUUID:       f.UserUUID,
		IngredientUUID: f.IngredientUUID,
		Amount:         f.Amount,
		Unit:           f.Unit,
		PurchasedDate:  f.PurchasedDate,
		ExpirationDate: f.ExpirationDate,
	}
}
