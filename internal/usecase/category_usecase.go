package usecase

import (
	"context"
	"fin-manager/internal/domain"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, new_category domain.NewCategoryData) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
}

type CategoryUseCase struct {
	CategoryRepository CategoryRepository
}

func NewCategoryUseCase(repo CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{CategoryRepository: repo}
}

func (c *CategoryUseCase) CreateCategory(ctx context.Context, new_category domain.NewCategoryData) (*domain.Category, error) {
	return c.CategoryRepository.CreateCategory(ctx, new_category)
}

func (c *CategoryUseCase) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return c.CategoryRepository.GetCategories(ctx)
}
