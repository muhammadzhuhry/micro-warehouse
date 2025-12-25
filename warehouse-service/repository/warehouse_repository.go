package repository

import (
	"context"
	"errors"
	"micro-warehouse/warehouse-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type WarehouseRepositoryInterface interface {
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	GetAllWarehouses(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Warehouse, int64, error)
	GetWarehouseByID(ctx context.Context, id uint) (*model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	DeleteWarehouse(ctx context.Context, id uint) error
}

type WarehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepositoryInterface {
	return &WarehouseRepository{db: db}
}

// CreateWarehouse implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] CreateWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return w.db.WithContext(ctx).Create(warehouse).Error
	}
}

// DeleteWarehouse implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) DeleteWarehouse(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] DeleteWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouse := &model.Warehouse{}

		// Check if warehouse exists
		if err := w.db.WithContext(ctx).Where("id = ?", id).Preload("WarehouseProducts").First(&modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] DeleteWarehouse - 2: %v", err)
			return err
		}

		if len(modelWarehouse.WarehouseProducts) > 0 {
			log.Errorf("[WarehouseRepository] DeleteWarehouse - 3: %v", errors.New("warehouse has product"))
			return errors.New("warehouse has product")
		}

		return w.db.WithContext(ctx).Delete(&model.Warehouse{}, id).Error
	}
}

// GetAllWarehouses implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) GetAllWarehouses(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Warehouse, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] GetAllWarehouses - 1: %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
		if page < 1 {
			page = 1
		}

		if limit < 1 {
			limit = 10
		}

		if sortBy == "" {
			sortBy = "created_at"
		}

		if sortOrder == "" {
			sortOrder = "desc"
		}

		// Calculate offset
		offset := (page - 1) * limit

		// Build query
		query := w.db.Model(&model.Warehouse{})

		// Add search filter if provided
		if search != "" {
			query = query.Where("name ILIKE ? OR address ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		var warehouses []model.Warehouse
		var totalRecords int64

		// Get total count
		if err := query.Count(&totalRecords).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetAllWarehouses - 2: %v", err)
			return nil, 0, err
		}

		if err := query.WithContext(ctx).
			Preload("WarehouseProducts").
			Order(sortBy + " " + sortOrder).
			Offset(offset).
			Limit(limit).
			Find(&warehouses).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetAllWarehouses - 3: %v", err)
			return nil, 0, err
		}
		return warehouses, totalRecords, nil
	}
}

// GetWarehouseByID implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) GetWarehouseByID(ctx context.Context, id uint) (*model.Warehouse, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] GetWarehouseByID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelWarehouse := model.Warehouse{}

		if err := w.db.WithContext(ctx).Where("id = ?", id).Preload("WarehouseProducts").First(&modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetWarehouseByID - 2: %v", err)
			return nil, err
		}
		return &modelWarehouse, nil
	}
}

// UpdateWarehouse implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] UpdateWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouse := &model.Warehouse{}

		// Check if warehouse exists
		if err := w.db.WithContext(ctx).Where("id = ?", warehouse.ID).First(&modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] UpdateWarehouse - 2: %v", err)
			return err
		}

		modelWarehouse.Name = warehouse.Name
		modelWarehouse.Address = warehouse.Address
		modelWarehouse.Phone = warehouse.Phone
		modelWarehouse.Photo = warehouse.Photo

		return w.db.WithContext(ctx).Save(&modelWarehouse).Error
	}
}
