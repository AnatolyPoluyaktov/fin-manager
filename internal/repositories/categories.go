package repositories

import (
	"context"
	"fin-manager/internal/domain"
	"fin-manager/internal/storage/pg"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	DB *pg.DB
}

func NewCategoryRepository(db *pg.DB) (*CategoryRepository, error) {
	return &CategoryRepository{DB: db}, nil
}

func (c *CategoryRepository) CreateCategory(ctx context.Context, new_category domain.NewCategoryData) (*domain.Category, error) {
	tx, err := c.DB.Pool.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("error rolling back transaction: %v", err)
		}
	}(tx, ctx)
	createdCategory, err := domain.NewCategory(new_category)
	if err != nil {
		return nil, fmt.Errorf("error creating new category: %v", err)
	}
	query := `INSERT INTO categories (name, description)
			          VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(ctx, query, new_category.Name, new_category.Description).Scan(&createdCategory.ID)

	if err != nil {
		return nil, fmt.Errorf("error inserting category: %v", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return &createdCategory, nil
}

func (c *CategoryRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	query := `SELECT id, name, description FROM categories`
	rows, err := c.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting categories: %v", err)
	}
	defer rows.Close()
	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err = rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, fmt.Errorf("error scanning category: %v", err)
		}
		categories = append(categories, category)

	}
	return categories, nil
}
