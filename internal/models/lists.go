package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type List struct {
	ID        int64
	UserID    int64
	StartDate time.Time
	Recipes   []*Recipe
}

// TODO: ALLOW USER TO ADD BLACKLIST FOR INGREDIENTS THEY DONT WANT ON THEIR LIST

func (l *List) GetIngredientsList() IngredientList {
	ingredients := make([]*Ingredient, 0)
	for _, r := range l.Recipes {
		ingredients = append(ingredients, r.Ingredients...)
	}

	return ingredients
}

type IngredientList []*Ingredient

func (l *IngredientList) GenerateMessage() (string, error) {
	var list strings.Builder
	_, err := list.WriteString("Here's your list:\n")
	if err != nil {
		return "", err
	}

	for _, i := range *l {
		_, err := list.WriteString(fmt.Sprintf("- %s\n", i.Name))
		if err != nil {
			return "", err
		}
	}

	_, err = list.WriteString("Send HELP for help")
	if err != nil {
		return "", err
	}

	return list.String(), nil
}

// func (l *IngredientsList) removeDuplicateIngredients()
// func (l *) removeExcludedIngredients()

type ListModel struct {
	db *sql.DB
}

func (m *ListModel) GetCurrentList(userID int64) (*List, error) {
	list := &List{UserID: userID}

	getList := `
	SELECT id, start_date
		FROM lists
		WHERE user_id = $1 AND current = TRUE`

	getRecipes := `
		SELECT id, user_id, title
		FROM recipes
		WHERE list_id = $1 AND user_id = $2`

	getIngredients := `
		SELECT id, quantity, unit, name
		FROM ingredients
		WHERE recipe_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, getList, userID).Scan(&list.ID, &list.StartDate)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, getRecipes, list.ID, userID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	recipes := make([]*Recipe, 0)
	for rows.Next() {
		r := new(Recipe)
		err := rows.Scan(&r.ID, &r.UserID, &r.Title)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		recipes = append(recipes, r)
	}
	list.Recipes = recipes
	rows.Close()

	for _, r := range list.Recipes {
		rows, err := tx.QueryContext(ctx, getIngredients, r.ID, r.UserID)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		ingredients := make([]*Ingredient, 0)
		for rows.Next() {
			i := new(Ingredient)
			err := rows.Scan(&i.ID, &i.Quantity, &i.Unit, &i.Name)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
			ingredients = append(ingredients, i)
		}
		r.Ingredients = ingredients
	}
	rows.Close()

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return list, nil
}

func (m *ListModel) StartNewList(userID int64) error {
	updateOldList := `
		UPDATE lists
		SET current = FALSE
		WHERE current = TRUE AND user_id = $1`

	insertNewList := `
		INSERT INTO lists (user_id)
		VALUES ($1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, updateOldList, userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, insertNewList, userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
