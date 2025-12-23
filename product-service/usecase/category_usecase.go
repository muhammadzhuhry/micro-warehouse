package usecase

import (
	"context"
	"micro-warehouse/product-service/model"
)

type CategoryUsecaseInterface interface {
	CreateCategory(ctx context.Context, category model.Category) error
	GetAllCategories(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Category, int64, error)
	GetCategoryByID(ctx context.Context, id uint) (*model.Category, error)
	UpdateCategory(ctx context.Context, category model.Category) error
	DeleteCategory(ctx context.Context, id uint) error
}

type CategoryUsecase struct {
	categoryRepo CategoryUsecaseInterface
}

func NewCategoryUsecase(categoryRepo CategoryUsecaseInterface) CategoryUsecaseInterface {
	return &CategoryUsecase{categoryRepo: categoryRepo}
}

// CreateCategory implements [CategoryRepositoryInterface].
func (c *CategoryUsecase) CreateCategory(ctx context.Context, category model.Category) error {
	return c.categoryRepo.CreateCategory(ctx, category)
}

// DeleteCategory implements [CategoryRepositoryInterface].
func (c *CategoryUsecase) DeleteCategory(ctx context.Context, id uint) error {
	return c.categoryRepo.DeleteCategory(ctx, id)
}

// GetAllCategories implements [CategoryRepositoryInterface].
func (c *CategoryUsecase) GetAllCategories(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Category, int64, error) {
	return c.categoryRepo.GetAllCategories(ctx, page, limit, search, sortBy, sortOrder)
}

// GetCategoryByID implements [CategoryRepositoryInterface].
func (c *CategoryUsecase) GetCategoryByID(ctx context.Context, id uint) (*model.Category, error) {
	return c.categoryRepo.GetCategoryByID(ctx, id)
}

// UpdateCategory implements [CategoryRepositoryInterface].
func (c *CategoryUsecase) UpdateCategory(ctx context.Context, category model.Category) error {
	return c.categoryRepo.UpdateCategory(ctx, category)
}
