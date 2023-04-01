package postgres

const sqlGetIngredient = `
	SELECT 	ingredient_uuid,
			ingredient_name,
			category,
			days_until_exp,
			created_at,
			updated_at
	
	FROM 	wdiet.ingredients
	
	WHERE	ingredient_uuid = $1

	LIMIT 1
	;
`

const sqlsearchIngredients = `
	SELECT 	ingredient_uuid,
			ingredient_name,
			category,
			days_until_exp,
			created_at,
			updated_at

	FROM wdiet.ingredients 
`

const sqlCreateIngredient = `
	INSERT INTO wdiet.ingredients(
		ingredient_name,
		category,
		days_until_exp
	)
	VALUES(
		$1,
		$2,
		$3
	)
	RETURNING ingredient_uuid, ingredient_name, category, days_until_exp, created_at, updated_at
	;
`

const sqlUpdateIngredient = `
	UPDATE wdiet.ingredients
		SET 
			ingredient_name = $1,
			category = $2,
			days_until_exp = $3,
			updated_at = now()
	WHERE ingredient_uuid = $4
	RETURNING ingredient_uuid, ingredient_name, category, days_until_exp, created_at, updated_at
	;
`

const sqlDeleteIngredient = `
	DELETE 
		FROM wdiet.ingredients

	WHERE ingredient_uuid = $1
	;
`

const sqlListFridgeIngredients = `
	SELECT 	user_uuid,
			ingredient_uuid,
			amount,
			unit,
			purchased_date,
			expiration_date,
			created_at,
			updated_at
	
	FROM 	wdiet.fridge_ingredients
	
	WHERE	user_uuid = $1
	;
`

const sqlCreateFridgeIngredient = `
	INSERT INTO wdiet.fridge_ingredients(
				user_uuid,
				ingredient_uuid,
				amount,
				unit,
				purchased_date,
				expiration_date
	)
	VALUES(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
	)
	RETURNING user_uuid, ingredient_uuid, amount, unit, purchased_date, expiration_date, created_at, updated_at
	;
`

const sqlUpdateFridgeIngredient = `
	UPDATE wdiet.fridge_ingredients
		SET 
			amount = $1,
			unit = $2,
			purchased_date = $3,
			updated_at = now()
	WHERE user_uuid = $4 AND ingredient_uuid = $5
	RETURNING user_uuid, ingredient_uuid, amount, unit, purchased_date, expiration_date, created_at, updated_at
	;
`

const sqlDeleteFridgeIngredient = `
	DELETE 
		FROM wdiet.fridge_ingredients

	WHERE user_uuid = $1 AND ingredient_uuid = $2
	;
`
