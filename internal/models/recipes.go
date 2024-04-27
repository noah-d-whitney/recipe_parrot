package models

import (
	"context"
	"database/sql"
	"time"
)

type Recipe struct {
	ID          int64
	UserID      int64
	Title       string
	Ingredients []Ingredient
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
		SELECT id, user_id, title
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

	err = tx.QueryRowContext(ctx, getRecipe, recipeID, userID).Scan(&recipe.ID, &recipe.UserID, &recipe.Title)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, getIngredients, recipeID, userID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	ingredients := make([]Ingredient, 0)
	for rows.Next() {
		i := Ingredient{}
		err := rows.Scan(&i.ID, &i.Quantity, &i.Unit, &i.Name)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		ingredients = append(ingredients, i)
	}
	recipe.Ingredients = ingredients

	return recipe, nil
}
func (m *RecipeModel) Create(r *Recipe) error {
	insertRecipe := `
		INSERT INTO recipes (title, user_id)
		VALUES ($1, $2)
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

	err = tx.QueryRowContext(ctx, insertRecipe, r.Title, r.UserID).Scan(&r.ID)
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
	}

	return nil
}