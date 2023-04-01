package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"wd_ingredients/store"

	"github.com/google/uuid"
)

const defaultTimeout = 5 * time.Second

func (pg *PG) GetIngredient(ctx context.Context, id uuid.UUID) (*store.Ingredient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var ingredient store.Ingredient

	row := pg.db.QueryRowContext(ctx, sqlGetIngredient, id)
	if err := row.Scan(
		&ingredient.IngredientUUID,
		&ingredient.IngredientName,
		&ingredient.Category,
		&ingredient.DaysUntilExp,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("error getting ingredient: %w", err)
	}

	return &ingredient, nil
}

func (pg *PG) SearchIngredients(ctx context.Context, i store.SearchIngredient) ([]store.Ingredient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var wheres []string
	var vars []interface{}
	var count int

	if i.IngredientName != nil {
		count++
		wheres = append(wheres, fmt.Sprintf(" ingredient_name = $%d", count))
		vars = append(vars, i.IngredientName)
	}
	if i.Category != nil {
		count++
		wheres = append(wheres, fmt.Sprintf(" category = $%d", count))
		vars = append(vars, i.Category)
	}

	whereClause := strings.Join(wheres, " AND ") //만약 하나가 안온다 그럼 join 자체가 안되기 때문에 AND로 묶이지도 않나보네..? 오호 go playground에서 테스트 해밨는데 맞음!!ㅋㅋ

	var ingredients []store.Ingredient

	fmt.Println("query is ", sqlsearchIngredients+" WHERE "+whereClause, "vars are ", vars)
	rows, err := pg.db.QueryContext(ctx, sqlsearchIngredients+" WHERE "+whereClause, vars...)
	if err != nil {
		return nil, fmt.Errorf("error searching ingredients: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ingredient store.Ingredient
		if err := rows.Scan(
			&ingredient.IngredientUUID,
			&ingredient.IngredientName,
			&ingredient.Category,
			&ingredient.DaysUntilExp,
			&ingredient.CreatedAt,
			&ingredient.UpdatedAt,
		); err != nil {
			// if errors.Is(err, sql.ErrNoRows) {
			// 	return nil, store.ErrNotFound
			// }
			return nil, fmt.Errorf("error searching ingredients: %w", err)
		}
		ingredients = append(ingredients, ingredient)
	}
	// if err := rows.Close(); err != nil {
	// 	return nil, err
	// }
	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }
	return ingredients, nil
}

func (pg *PG) CreateIngredient(ctx context.Context, i store.Ingredient) (*store.Ingredient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating ingredient: %w", err)
	}

	var ingredient store.Ingredient

	row := tx.QueryRowContext(ctx, sqlCreateIngredient,
		&i.IngredientName,
		&i.Category,
		&i.DaysUntilExp,
	)

	if err = row.Scan(
		&ingredient.IngredientUUID,
		&ingredient.IngredientName,
		&ingredient.Category,
		&ingredient.DaysUntilExp,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
	); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating ingredient: %w", err)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating ingredient: %w", err)
	}

	return &ingredient, nil
}

func (pg *PG) UpdateIngredient(ctx context.Context, i store.Ingredient) (*store.Ingredient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error updating ingredient: %w", err)
	}

	var ingredient store.Ingredient

	row := tx.QueryRowContext(ctx, sqlUpdateIngredient,
		i.IngredientName,
		i.Category,
		i.DaysUntilExp,
		i.IngredientUUID,
	)

	if err = row.Scan(
		&ingredient.IngredientUUID,
		&ingredient.IngredientName,
		&ingredient.Category,
		&ingredient.DaysUntilExp,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil, store.ErrNotFound
		}
		tx.Rollback()
		return nil, fmt.Errorf("error updating ingredient: %w", err)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating ingredient: %w", err)
	}

	return &ingredient, nil
}

func (pg *PG) DeleteIngredient(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel() //it defers cancel until it returns from an error(actually, until it returns from current function, even john doesn't know what will happen after you return without an error. Does cancel() execute?) So what does it do? It cancels the context everywhere, so you don't do extra work when you know it's gonna fail. Context is per request.

	//when writing to or deleting from db, you should be in a transaction
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error deleting ingredient: %w", err)
	}

	res, err := tx.ExecContext(ctx, sqlDeleteIngredient, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting ingredient: %w", err)
	}

	if affected, _ := res.RowsAffected(); affected != 1 { //you do this to prevent update when without WHERE clause, or if there's two rows that have the same WHERE condition. I don't need to do it when it's queryRowContext. Only works with multiple rows.
		tx.Rollback()
		return fmt.Errorf("error deleting ingredient, rows affected is %d instead of 1", affected)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting ingredient: %w", err)
	}

	return nil
}

func (pg *PG) ListFridgeIngredients(ctx context.Context, id uuid.UUID) ([]store.FridgeIngredient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var fridgeIngredients []store.FridgeIngredient

	rows, err := pg.db.QueryContext(ctx, sqlListFridgeIngredients, id)
	if err != nil {
		return nil, fmt.Errorf("error listing fridge ingredients: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fridgeIngredient store.FridgeIngredient
		if err := rows.Scan(
			&fridgeIngredient.UserUUID,
			&fridgeIngredient.IngredientUUID,
			&fridgeIngredient.Amount,
			&fridgeIngredient.Unit,
			&fridgeIngredient.PurchasedDate,
			&fridgeIngredient.ExpirationDate,
			&fridgeIngredient.CreatedAt,
			&fridgeIngredient.UpdatedAt,
		); err != nil {
			// if errors.Is(err, sql.ErrNoRows) {
			// 	return nil, store.ErrNotFound
			// }
			return nil, fmt.Errorf("error listing fridge ingredients: %w", err)
		}
		fridgeIngredients = append(fridgeIngredients, fridgeIngredient)
	}

	return fridgeIngredients, nil
}

func (pg *PG) CreateFridgeIngredient(ctx context.Context, f store.FridgeIngredient) (*store.FridgeIngredient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating fridge ingredient: %w", err)
	}

	var fridgeIngredient store.FridgeIngredient

	row := tx.QueryRowContext(ctx, sqlCreateFridgeIngredient,
		&f.UserUUID,
		&f.IngredientUUID,
		&f.Amount,
		&f.Unit,
		&f.PurchasedDate,
		&f.ExpirationDate,
	)

	err = row.Scan(
		&fridgeIngredient.UserUUID,
		&fridgeIngredient.IngredientUUID,
		&fridgeIngredient.Amount,
		&fridgeIngredient.Unit,
		&fridgeIngredient.PurchasedDate,
		&fridgeIngredient.ExpirationDate,
		&fridgeIngredient.CreatedAt,
		&fridgeIngredient.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating fridge ingredient: %w", err)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating fridge ingredient: %w", err)
	}

	return &fridgeIngredient, nil
}

func (pg *PG) UpdateFridgeIngredient(ctx context.Context, f store.FridgeIngredient) (*store.FridgeIngredient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error updating fridge ingredient: %w", err)
	}

	var fridgeIngredient store.FridgeIngredient

	row := tx.QueryRowContext(ctx, sqlUpdateFridgeIngredient, //여기에 나온 순서는 sql.go에서 $로 순서매겨놓은 순서.
		&f.Amount,
		&f.Unit,
		&f.PurchasedDate,
		&f.UserUUID,
		&f.IngredientUUID,
	)

	if err = row.Scan( //여기에 나온 순서는 sql.go에서 RETURNING 뒤의 필드 순서.
		&fridgeIngredient.UserUUID,
		&fridgeIngredient.IngredientUUID,
		&fridgeIngredient.Amount,
		&fridgeIngredient.Unit,
		&fridgeIngredient.PurchasedDate,
		&fridgeIngredient.ExpirationDate, //update fridge ingredient 핸들러에서 f.purchased date로 산출된 expiration date를 여기에 넣어줘야함.. 헐 바보네 &f.ExpirationDate가 아니지..-_- f.ExpirationDate는 이미 핸들러에서 바뀌어서 request body에 있었으므로 row 안에 이미 포함됐고, 그 내용을 &fridgeIngredient.ExpirationDate에 붙여넣어줘야지..
		&fridgeIngredient.CreatedAt,
		&fridgeIngredient.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil, store.ErrNotFound
		}
		tx.Rollback()
		return nil, fmt.Errorf("error updating fridge ingredient: %w", err)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating fridge ingredient: %w", err)
	}

	return &fridgeIngredient, nil
}

func (pg *PG) DeleteFridgeIngredient(ctx context.Context, uid uuid.UUID, fid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel() //to make sure that the cancel function runs, otherwise I have to say it before every return caused by error

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error deleting fridge ingredient: %w", err)
	}

	res, err := tx.ExecContext(ctx, sqlDeleteFridgeIngredient, uid, fid)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting fridge ingredient: %w", err)
	}

	if affected, _ := res.RowsAffected(); affected != 1 {
		tx.Rollback()
		return fmt.Errorf("error deleting fridge ingredient, rows affected %d instead of 1", affected)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting fridge ingredient: %w", err)
	}

	return nil
}
