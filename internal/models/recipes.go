package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Recipe struct {
	ID          int64
	UserID      int64
	ListID      int64
	URL         string
	Title       string
	Ingredients []*Ingredient
}

type Ingredient struct {
	ID       int64
	Quantity string
	Unit     string
	Name     string
}

type RecipeModel struct {
	db *sql.DB
}

func (m *RecipeModel) Get(recipeID, userID int64) (*Recipe, error) {
	getRecipe := `
		SELECT id, user_id, title, list_id, url
		FROM recipes
		WHERE id = $1 AND user_id = $2`

	getIngredients := `
		SELECT id, quantity, unit, name
		FROM ingredients
		WHERE recipe_id = $1 AND user_id = $2`

	recipe := new(Recipe)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, getRecipe, recipeID, userID).Scan(&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.ListID, &recipe.URL)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, getIngredients, recipeID, userID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	ingredients := make([]*Ingredient, 0)
	for rows.Next() {
		i := &Ingredient{}
		err := rows.Scan(&i.ID, &i.Quantity, &i.Unit, &i.Name)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		ingredients = append(ingredients, i)
	}
	recipe.Ingredients = ingredients

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return recipe, nil
}

func (m *RecipeModel) Create(r *Recipe) error {
	getCurrentListID := `SELECT id FROM lists WHERE current = TRUE AND user_id = $1`

	insertRecipe := `
		INSERT INTO recipes (title, user_id, list_id, url)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	insertIngredient := `
		INSERT INTO ingredients (user_id, recipe_id, quantity, unit, name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, getCurrentListID, r.UserID).Scan(&r.ListID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.QueryRowContext(ctx, insertRecipe, r.Title, r.UserID, r.ListID, r.URL).Scan(&r.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, i := range r.Ingredients {
		args := []any{r.UserID, r.ID, i.Quantity, i.Unit, i.Name}
		err := tx.QueryRowContext(ctx, insertIngredient, args...).Scan(&i.ID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		fmt.Printf("INGREDIENT ID: %d\n", i.ID)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
