package repository

import (
	"context"
	"micro-warehouse/warehouse-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// Detail warehouse, Get warehouse product by warehouseID and productID, Get Product total Stock
// Create warehouse product, Update warehouse product, Delete warehouse product
// Delete All Product in Warehouse Product

type WarehouseProductRepositoryInterface interface {
	CreateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error
	GetWarehouseProductByWarehouseIDAndProductID(ctx context.Context, warehouseID, productID uint) (*model.WarehouseProduct, error)
	UpdateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error
	DeleteWarehouseProduct(ctx context.Context, warehouseProductID uint) error
	DeleteAllWarehouseProductByProductID(ctx context.Context, productID uint) error
	GetWarehouseProductByProductID(ctx context.Context, productID uint) ([]model.WarehouseProduct, error)
	GetProductTotalStock(ctx context.Context, productID uint) (int, error)
	GetDetailWarehouse(ctx context.Context, warehouseID uint) (*model.Warehouse, error)
	GetDetailWarehouseProductByID(ctx context.Context, warehouseID uint) (*model.WarehouseProduct, error)
}

type WarehouseProductRepository struct {
	db *gorm.DB
}

func NewWarehouseProductRepository(db *gorm.DB) WarehouseProductRepositoryInterface {
	return &WarehouseProductRepository{db: db}
}

// CreateWarehouseProduct implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) CreateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] CreateWarehouseProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return w.db.WithContext(ctx).Create(warehouseProduct).Error
	}
}

// DeleteAllWarehouseProductByProductID implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) DeleteAllWarehouseProductByProductID(ctx context.Context, productID uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] DeleteAllWarehouseProductByProductID - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		err := w.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&model.WarehouseProduct{}).Error
		if err != nil {
			log.Errorf("[WarehouseProductRepository] DeleteAllWarehouseProductByProductID - 2: %v", err)
			return err
		}

		return nil
	}
}

// DeleteWarehouseProduct implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) DeleteWarehouseProduct(ctx context.Context, warehouseProductID uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] DeleteWarehouseProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouseProduct := &model.WarehouseProduct{}

		// Check if warehouse product exists
		if err := w.db.WithContext(ctx).Where("id = ?", warehouseProductID).First(&modelWarehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] DeleteWarehouseProduct - 2: %v", err)
			return err
		}

		return w.db.WithContext(ctx).Delete(&model.WarehouseProduct{}).Error
	}
}

// GetDetailWarehouse implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetDetailWarehouse(ctx context.Context, warehouseID uint) (*model.Warehouse, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetDetailWarehouse - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var warehouse model.Warehouse

		if err := w.db.WithContext(ctx).Where("id = ?", warehouseID).Select("id", "name", "photo", "phone").
			Preload("WarehouseProducts").
			First(&warehouse).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetDetailWarehouse - 2: %v", err)
			return nil, err
		}

		return &warehouse, nil
	}
}

// GetDetailWarehouseProductByID implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetDetailWarehouseProductByID(ctx context.Context, warehouseID uint) (*model.WarehouseProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetDetailWarehouseProductByID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var modelWarehouseProduct model.WarehouseProduct
		if err := w.db.WithContext(ctx).Where("warehouse_id = ?", warehouseID).
			Select("id", "product_id", "stock", "warehouse_id").
			Preload("Warehouse").
			First(&modelWarehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetDetailWarehouseProductByID - 2: %v", err)
			return nil, err
		}

		return &modelWarehouseProduct, nil
	}
}

// GetProductTotalStock implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetProductTotalStock(ctx context.Context, productID uint) (int, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetProductTotalStock - 1: %v", ctx.Err())
		return 0, ctx.Err()
	default:
		var totalStock int
		if err := w.db.WithContext(ctx).
			Model(&model.WarehouseProduct{}).
			Where("product_id = ?", productID).
			Select("COALESCE(SUM(stock), 0)").
			Scan(&totalStock).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetProductTotalStock - 2: %v", err)
			return 0, err
		}

		return totalStock, nil
	}
}

// GetWarehouseProductByProductID implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetWarehouseProductByProductID(ctx context.Context, productID uint) ([]model.WarehouseProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetWarehouseProductByProductID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var warehouseProducts []model.WarehouseProduct

		if err := w.db.WithContext(ctx).
			Where("product_id = ?", productID).
			Preload("Warehouse").
			Find(&warehouseProducts).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetWarehouseProductByProductID - 2: %v", err)
			return nil, err
		}

		return warehouseProducts, nil
	}
}

// GetWarehouseProductByWarehouseIDAndProductID implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetWarehouseProductByWarehouseIDAndProductID(ctx context.Context, warehouseID uint, productID uint) (*model.WarehouseProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetWarehouseProductByWarehouseIDAndProductID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var warehouseProduct model.WarehouseProduct

		if err := w.db.WithContext(ctx).
			Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
			First(&warehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetWarehouseProductByWarehouseIDAndProductID - 2: %v", err)
			return nil, err
		}

		return &warehouseProduct, nil
	}
}

// UpdateWarehouseProduct implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) UpdateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] UpdateWarehouseProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouseProduct := &model.WarehouseProduct{}

		// Check if warehouse product exists
		if err := w.db.WithContext(ctx).Where("id = ?", warehouseProduct.ID).First(&modelWarehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] UpdateWarehouseProduct - 2: %v", err)
			return err
		}

		modelWarehouseProduct.Stock = warehouseProduct.Stock
		modelWarehouseProduct.WarehouseID = warehouseProduct.WarehouseID
		modelWarehouseProduct.ProductID = warehouseProduct.ProductID

		return w.db.WithContext(ctx).Save(warehouseProduct).Error
	}
}
