package recipe

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

func NewRepo(db *sql.DB) *RepoSQL {
	return &RepoSQL{
		db: db,
	}
}

type RepoSQL struct {
	db *sql.DB
}

func (r *RepoSQL) FindByID(ctx context.Context, id int64) (*entity.Recipe, error) {
	// Get recipe steps
	rows, err := r.db.QueryContext(ctx, `SELECT content FROM recipeSteps WHERE recipeId = ? ORDER BY stepNo ASC`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	var steps []string
	for i := 0; rows.Next(); i++ {
		var temp string
		rows.Scan(&temp)

		steps = append(steps, temp)
	}

	rows, err = r.db.QueryContext(ctx, `
		SELECT
			ri.ingredientId, i.name, i.type
		FROM recipeIngredients ri
		JOIN ingredients i ON (ri.ingredientId = i.id)
		WHERE ri.recipeId = ?`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	var ingredients []*entity.Ingredient
	for rows.Next() {
		ingredient := &entity.Ingredient{}
		rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Type)
		ingredients = append(ingredients, ingredient)
	}

	// Get recipe
	row := r.db.QueryRowContext(ctx, `SELECT id, name FROM recipes WHERE id = ?`, id)
	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	recipe := &entity.Recipe{
		Ingredients: ingredients,
		Steps:       steps,
	}

	row.Scan(&recipe.ID, &recipe.Name)
	return recipe, nil
}

func (r *RepoSQL) Add(ctx context.Context, recipe *entity.Recipe) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `INSERT INTO recipes (name) VALUES (?)`, recipe.Name)
	if err != nil {
		return err
	}

	recipe.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	for i, step := range recipe.Steps {
		tx.ExecContext(ctx, `INSERT INTO recipeSteps (recipeId, stepNo, content) VALUES (?, ?, ?)`,
			recipe.ID, i, step)
	}

	for _, ingredient := range recipe.Ingredients {
		tx.ExecContext(ctx, `INSERT INTO recipeIngredients (recipeId, ingredientId) VALUES (?, ?)`,
			recipe.ID, ingredient.ID)
	}

	return tx.Commit()
}
