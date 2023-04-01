package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"wd_ingredients/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) GetIngredient(c *gin.Context) {
	l := s.l.Named("GetIngredient")

	id := c.Param("id")

	iid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error getting ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	ingredient, err := s.db.GetIngredient(context.Background(), iid)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			l.Info("error getting ingredient", zap.Error(err))
			c.Status(http.StatusNotFound)
			return
		}
		l.Error("error getting ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbIngr2ApiIngr(ingredient))
}

func (s *Service) SearchIngredients(c *gin.Context) {
	l := s.l.Named("SearchIngredients")

	var searchIngrRequest SearchIngredients

	if err := json.NewDecoder(c.Request.Body).Decode(&searchIngrRequest); err != nil {
		l.Info("error searching ingredients", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidSearchIngrRequest(searchIngrRequest) {
		l.Info("error searching ingredients")
		c.Status(http.StatusBadRequest)
		return
	}

	ingredients, err := s.db.SearchIngredients(context.Background(), apiSearchIngr2DBSearchIngr(searchIngrRequest))
	if err != nil {
		// if errors.Is(err, store.ErrNotFound) {
		// 	l.Info("error searching ingredients", zap.Error(err))
		// 	c.Status(http.StatusNotFound)
		// 	return
		// }
		l.Error("error searching ingredients", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(ingredients) == 0 {
		c.Status(http.StatusOK)
		return
	}

	var searchIngrResponse []Ingredient

	for _, ingredient := range ingredients {
		i := dbIngr2ApiIngr(&ingredient)
		searchIngrResponse = append(searchIngrResponse, i)
	}

	c.JSON(http.StatusOK, searchIngrResponse)
}

func (s *Service) CreateIngredient(c *gin.Context) {
	l := s.l.Named("CreateIngredient")

	var createIngrRequest Ingredient

	if err := json.NewDecoder(c.Request.Body).Decode(&createIngrRequest); err != nil {
		l.Info("error creating ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidCreateIngrRequest(createIngrRequest) {
		l.Info("error creating ingredient")
		c.Status(http.StatusBadRequest)
		return
	}

	ingredient, err := s.db.CreateIngredient(context.Background(), apiIngr2DBIngr(createIngrRequest))
	if err != nil {
		l.Error("error creating ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbIngr2ApiIngr(ingredient))
}

func (s *Service) UpdateIngredient(c *gin.Context) {
	l := s.l.Named("UpdateIngredient")

	id := c.Param("id")

	iid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error updating ingredient", zap.Error(err)) //error message shouldn't contain single quote(') cause it might break. Spacebar is okay.
		c.Status(http.StatusBadRequest)
		return
	}

	var updateIngrRequest Ingredient

	if err := json.NewDecoder(c.Request.Body).Decode(&updateIngrRequest); err != nil {
		l.Info("error updating ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidUpdateIngrRequest(updateIngrRequest, iid) {
		l.Info("error updating ingredient")
		c.Status(http.StatusBadRequest)
		return
	}

	ingredient, err := s.db.UpdateIngredient(context.Background(), apiIngr2DBIngr(updateIngrRequest)) //if I have two variables, I can still do combined if statement, like if user, err ... ; err != nil {}, but then user can only survive within the next 3 lines of if statement. So I can't return user variable at the bottom in c.JSON().
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			l.Info("error updating ingredient", zap.Error(err))
			c.Status(http.StatusNotFound)
			return
		}
		l.Error("error updating ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbIngr2ApiIngr(ingredient))
}

func (s *Service) DeleteIngredient(c *gin.Context) {
	l := s.l.Named("DeleteIngredient")

	id := c.Param("id")

	iid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error deleting ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if err = s.db.DeleteIngredient(context.Background(), iid); err != nil {
		l.Error("error deleting ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Service) ListFridgeIngredients(c *gin.Context) {
	l := s.l.Named("ListFridgeIngredients")

	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error listing fridge ingredients", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	fridgeIngredients, err := s.db.ListFridgeIngredients(context.Background(), uid)
	if err != nil {
		// if errors.Is(err, store.ErrNotFound) {
		// 	l.Info("error listing fridge ingredients", zap.Error(err))
		// 	c.Status(http.StatusNotFound)
		// 	return
		// }
		l.Error("error listing fridge ingredients", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(fridgeIngredients) == 0 {
		c.Status(http.StatusOK)
		return
	}

	var listFIngrResponse []FridgeIngredient

	for _, f := range fridgeIngredients {
		fridgeIngredient := dbFIngr2ApiFIngr(&f)
		listFIngrResponse = append(listFIngrResponse, fridgeIngredient)
	}

	c.JSON(http.StatusOK, listFIngrResponse)
}

func (s *Service) CreateFridgeIngredient(c *gin.Context) {
	l := s.l.Named("CreateFridge")

	var createFIngrRequest FridgeIngredient

	if err := json.NewDecoder(c.Request.Body).Decode(&createFIngrRequest); err != nil {
		l.Info("error creating fridge ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidCreateFIngrRequest(createFIngrRequest) {
		l.Info("error creating fridge ingredient")
		c.Status(http.StatusBadRequest)
		return
	}

	ingredient, err := s.db.GetIngredient(context.Background(), createFIngrRequest.IngredientUUID)
	if err != nil {
		l.Error("error creating fridge ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	createFIngrRequest.ExpirationDate = createFIngrRequest.PurchasedDate.Add(24 * time.Hour * time.Duration(ingredient.DaysUntilExp))

	fridgeIngredient, err := s.db.CreateFridgeIngredient(context.Background(), apiFIngr2DBFIngr(createFIngrRequest))
	if err != nil {
		l.Error("error creating fridge ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbFIngr2ApiFIngr(fridgeIngredient))
}

func (s *Service) UpdateFridgeIngredient(c *gin.Context) {
	l := s.l.Named("UpdateFridgeIngredient")

	id := c.Param("id")

	iid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error updating fridge ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	var updateFIngrRequest FridgeIngredient

	if err := json.NewDecoder(c.Request.Body).Decode(&updateFIngrRequest); err != nil {
		l.Info("error updating fridge ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidUpdateFIngrRequest(updateFIngrRequest, iid) {
		l.Info("error updating fridge ingredient")
		c.Status(http.StatusBadRequest)
		return
	}

	ingredient, err := s.db.GetIngredient(context.Background(), updateFIngrRequest.IngredientUUID)
	if err != nil {
		l.Error("error updating fridge ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	updateFIngrRequest.ExpirationDate = updateFIngrRequest.PurchasedDate.Add(24 * time.Hour * time.Duration(ingredient.DaysUntilExp))

	fridgeIngredient, err := s.db.UpdateFridgeIngredient(context.Background(), apiFIngr2DBFIngr(updateFIngrRequest))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			l.Info("error updating fridge ingredient", zap.Error(err))
			c.Status(http.StatusNotFound)
			return
		}
		l.Error("error updating fridge ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbFIngr2ApiFIngr(fridgeIngredient))
}

func (s *Service) DeleteFridgeIngredient(c *gin.Context) {
	l := s.l.Named("DeleteFridgeIngredient")

	id := c.Param("uid")
	id2 := c.Param("iid")

	uid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error deleting fridge ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}
	iid, err := uuid.Parse(id2)
	if err != nil {
		l.Info("error deleting fridge ingredient", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	// var deleteFIngrRequest DeleteFIngr

	// if err := json.NewDecoder(c.Request.Body).Decode(&deleteFIngrRequest); err != nil {
	// 	l.Info("error deleting fridge ingredient", zap.Error(err))
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// }

	// if !isValidDeleteFIngrRequest(deleteFIngrRequest, uid) {
	// 	l.Info("error deleting fridge ingredient")
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// }

	if err := s.db.DeleteFridgeIngredient(context.Background(), uid, iid); err != nil {
		l.Error("error deleting fridge ingredient", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
