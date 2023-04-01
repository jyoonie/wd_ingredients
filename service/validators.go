package service

import (
	"github.com/google/uuid"
)

func isValidSearchIngrRequest(i SearchIngredients) bool {
	if i.IngredientName == "" && i.Category == "" {
		return false
	}

	return true
}

func isValidCreateIngrRequest(i Ingredient) bool {
	switch {
	case i.IngredientUUID != uuid.Nil:
		return false
	case i.IngredientName == "":
		return false
	case !isCategory(i.Category):
		return false
	case i.DaysUntilExp < 0:
		return false
	}

	return true
}

func isCategory(c string) bool {
	switch {
	case c == "vegetables":
		return true
	case c == "fruits":
		return true
	case c == "meat":
		return true
	case c == "fish":
		return true
	case c == "eggs":
		return true
	case c == "dairy":
		return true
	case c == "grains":
		return true
	case c == "water":
		return true
	case c == "etc":
		return true
	}

	return false
}

func isValidUpdateIngrRequest(i Ingredient, uidFromPath uuid.UUID) bool {
	switch {
	case uidFromPath != i.IngredientUUID:
		return false
	case i.IngredientUUID == uuid.Nil:
		return false
	case i.IngredientName == "":
		return false
	case !isCategory(i.Category):
		return false
	case i.DaysUntilExp < 0:
		return false
	}

	return true
}

func isValidCreateFIngrRequest(f FridgeIngredient) bool {
	switch {
	case f.UserUUID == uuid.Nil:
		return false
	case f.IngredientUUID == uuid.Nil:
		return false
	case f.Amount <= 0:
		return false
	case f.Unit == "":
		return false
	case f.PurchasedDate.IsZero():
		return false
	case !f.ExpirationDate.IsZero():
		return false
	}

	return true
}

func isValidUpdateFIngrRequest(f FridgeIngredient, uidFromPath uuid.UUID) bool {
	switch {
	case uidFromPath != f.IngredientUUID:
		return false
	case f.UserUUID == uuid.Nil:
		return false
	case f.IngredientUUID == uuid.Nil:
		return false
	case f.Amount <= 0:
		return false
	case f.Unit == "":
		return false
	case f.PurchasedDate.IsZero():
		return false
	case !f.ExpirationDate.IsZero(): //always validate the data like the front end is retarded
		return false
	}

	return true
}
