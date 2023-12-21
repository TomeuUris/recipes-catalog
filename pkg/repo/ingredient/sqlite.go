package ingredient

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

func (r *RepoSQL) FindByID(ctx context.Context, id int) (*entity.Ingredient, error) {
	// Get ingredient from database
	row := r.db.QueryRowContext(ctx, `SELECT id, name, type FROM ingredients WHERE id = ?`, id)
	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	ingredient := &entity.Ingredient{}
	row.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Type)

	return ingredient, nil
}

func (r *RepoSQL) FindByFilter(ctx context.Context, f *FindFilter) ([]*entity.Ingredient, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM ingredients WHERE type = ?`, f.Type)
	if err != nil {
		return nil, err
	}

	var ingredients []*entity.Ingredient
	for rows.Next() {
		var ingredient *entity.Ingredient
		rows.Scan(ingredient)
		ingredients = append(ingredients, &entity.Ingredient{
			ID:   ingredient.ID,
			Name: ingredient.Name,
			Type: ingredient.Type,
		})
	}

	return ingredients, nil
}

func (r *RepoSQL) Add(ctx context.Context, ingredient *entity.Ingredient) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `INSERT INTO ingredients (name, type) VALUES (?, ?)`, ingredient.Name, ingredient.Type)
	if err != nil {
		return err
	}

	ingredient.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return tx.Commit()
}
