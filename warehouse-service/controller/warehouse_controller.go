package controller

import (
	"micro-warehouse/warehouse-service/controller/request"
	"micro-warehouse/warehouse-service/controller/response"
	"micro-warehouse/warehouse-service/model"
	"micro-warehouse/warehouse-service/pkg/conv"
	"micro-warehouse/warehouse-service/pkg/pagination"
	"micro-warehouse/warehouse-service/pkg/validator"
	"micro-warehouse/warehouse-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type WarehouseControllerInterface interface {
	CreateWarehouse(c *fiber.Ctx) error
	GetAllWarehouses(c *fiber.Ctx) error
	GetWarehouseByID(c *fiber.Ctx) error
	UpdateWarehouse(c *fiber.Ctx) error
	DeleteWarehouse(c *fiber.Ctx) error
}

type WarehouseController struct {
	warehouseUsecase usecase.WarehouseUsecaseInterface
}

func NewWarehouseController(warehouseUsecase usecase.WarehouseUsecaseInterface) WarehouseControllerInterface {
	return &WarehouseController{
		warehouseUsecase: warehouseUsecase,
	}
}

// CreateWarehouse implements [WarehouseControllerInterface].
func (w *WarehouseController) CreateWarehouse(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.CreateWarehouseRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	warehouseModel := model.Warehouse{
		Name:    req.Name,
		Address: req.Address,
		Photo:   req.Photo,
		Phone:   req.Phone,
	}

	if err := w.warehouseUsecase.CreateWarehouse(ctx, &warehouseModel); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Warehouse created successfully",
	})
}

// DeleteWarehouse implements [WarehouseControllerInterface].
func (w *WarehouseController) DeleteWarehouse(c *fiber.Ctx) error {
	ctx := c.Context()

	warehouseID := c.Params("id")
	if warehouseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "warehouse ID is required",
		})
	}

	id := conv.StringToUint(warehouseID)

	if err := w.warehouseUsecase.DeleteWarehouse(ctx, id); err != nil {
		log.Errorf("[WarehouseController] DeleteWarehouse - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete warehouse",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse deleted successfully",
	})
}

// GetAllWarehouses implements [WarehouseControllerInterface].
func (w *WarehouseController) GetAllWarehouses(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.GetAllWarehousesRequest{}
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouses - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouses - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	warehouses, totalRecords, err := w.warehouseUsecase.GetAllWarehouses(ctx, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouses - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get warehouses",
		})
	}

	resp := []response.WarehouseResponse{}
	for _, warehouse := range warehouses {
		resp = append(resp, response.WarehouseResponse{
			ID:           warehouse.ID,
			Name:         warehouse.Name,
			Address:      warehouse.Address,
			Photo:        warehouse.Photo,
			Phone:        warehouse.Phone,
			CountProduct: len(warehouse.WarehouseProducts),
		})
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(totalRecords))

	response := response.GetAllWarehousesResponse{
		Warehouses: resp,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Warehouses retrieved successfully",
	})
}

// GetWarehouseByID implements [WarehouseControllerInterface].
func (w *WarehouseController) GetWarehouseByID(c *fiber.Ctx) error {
	ctx := c.Context()

	warehouseID := c.Params("id")
	if warehouseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "warehouse ID is required",
		})
	}

	id := conv.StringToUint(warehouseID)

	warehouse, err := w.warehouseUsecase.GetWarehouseByID(ctx, id)
	if err != nil {
		log.Errorf("[WarehouseController] GetWarehouseByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get warehouse",
		})
	}

	resp := response.DetailWarehouseResponse{
		ID:      warehouse.ID,
		Name:    warehouse.Name,
		Address: warehouse.Address,
		Photo:   warehouse.Photo,
		Phone:   warehouse.Phone,
	}

	warehouseProductResp := []response.WarehouseProductResponse{}
	for _, wp := range warehouse.WarehouseProducts {
		warehouseProductResp = append(warehouseProductResp, response.WarehouseProductResponse{
			ID:          wp.ID,
			ProductID:   wp.ProductID,
			WarehouseID: wp.WarehouseID,
			Stock:       wp.Stock,
		})
	}

	resp.WarehouseProducts = warehouseProductResp

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Warehouse retrieved successfully",
	})
}

// UpdateWarehouse implements [WarehouseControllerInterface].
func (w *WarehouseController) UpdateWarehouse(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.CreateWarehouseRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	warehouseModel := model.Warehouse{
		ID:      conv.StringToUint(c.Params("id")),
		Name:    req.Name,
		Address: req.Address,
		Photo:   req.Photo,
		Phone:   req.Phone,
	}

	if err := w.warehouseUsecase.UpdateWarehouse(ctx, &warehouseModel); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse updated successfully",
	})
}
