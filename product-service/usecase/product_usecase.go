package usecase

import (
	"context"
	"micro-warehouse/product-service/model"
)

type ProductUsecaseInterface interface {
	CreateProduct(ctx context.Context, product model.Product) error
	GetAllProducts(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Product, int64, error)
	GetProductByID(ctx context.Context, id uint) (*model.Product, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*model.Product, error)
	UpdateProduct(ctx context.Context, product model.Product) error
	DeleteProduct(ctx context.Context, id uint) error
}

type ProductUsecase struct {
	productRepo ProductUsecaseInterface
}

func NewProductUsecase(productRepo ProductUsecaseInterface) ProductUsecaseInterface {
	return &ProductUsecase{productRepo: productRepo}
}

// CreateProduct implements [ProductUsecaseInterface].
func (p *ProductUsecase) CreateProduct(ctx context.Context, product model.Product) error {
	return p.productRepo.CreateProduct(ctx, product)
}

// DeleteProduct implements [ProductUsecaseInterface].
func (p *ProductUsecase) DeleteProduct(ctx context.Context, id uint) error {
	return p.productRepo.DeleteProduct(ctx, id)
}

// GetAllProducts implements [ProductUsecaseInterface].
func (p *ProductUsecase) GetAllProducts(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Product, int64, error) {
	return p.productRepo.GetAllProducts(ctx, page, limit, search, sortBy, sortOrder)
}

// GetProductByBarcode implements [ProductUsecaseInterface].
func (p *ProductUsecase) GetProductByBarcode(ctx context.Context, barcode string) (*model.Product, error) {
	return p.productRepo.GetProductByBarcode(ctx, barcode)
}

// GetProductByID implements [ProductUsecaseInterface].
func (p *ProductUsecase) GetProductByID(ctx context.Context, id uint) (*model.Product, error) {
	return p.productRepo.GetProductByID(ctx, id)
}

// UpdateProduct implements [ProductUsecaseInterface].
func (p *ProductUsecase) UpdateProduct(ctx context.Context, product model.Product) error {
	return p.productRepo.UpdateProduct(ctx, product)
}
