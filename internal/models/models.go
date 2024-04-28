package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users   UserModel
	Sites   SiteModel
	Recipes RecipeModel
	Lists   ListModel
}

func NewModels(initDB *sql.DB) Models {
	return Models{
		Users:   UserModel{db: initDB},
		Sites:   SiteModel{db: initDB},
		Recipes: RecipeModel{db: initDB},
		Lists:   ListModel{db: initDB},
	}
}
