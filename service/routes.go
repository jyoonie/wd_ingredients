package service

func (s *Service) registerRoutes() {
	s.r.POST("/ingredients/search", s.SearchIngredients)

	authorized := s.r.Group("/")
	//authorized.Use(s.ValidateToken)
	{
		authorized.GET("/ingredients/:id", s.GetIngredient)
		authorized.POST("/ingredients", s.CreateIngredient)
		authorized.POST("/ingredients/:id", s.UpdateIngredient)
		authorized.DELETE("/ingredients/:id", s.DeleteIngredient)

		authorized.GET("/users/:id/fridge_ingredients", s.ListFridgeIngredients)
		authorized.POST("/fridge_ingredients", s.CreateFridgeIngredient)
		authorized.POST("/fridge_ingredients/:id", s.UpdateFridgeIngredient)
		authorized.DELETE("/users/:uid/fridge_ingredients/:iid", s.DeleteFridgeIngredient)

	}
}
