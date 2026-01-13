package controller

import (
	"micro-warehouse/warehouse-service/controller/request"
	"micro-warehouse/warehouse-service/controller/response"
	"micro-warehouse/warehouse-service/model"
	"micro-warehouse/warehouse-service/pkg/conv"
	"micro-warehouse/warehouse-service/pkg/httpclient"
	"micro-warehouse/warehouse-service/pkg/validator"
	"micro-warehouse/warehouse-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type WarehouseProductControllerInterface interface {
	CreateWarehouseProduct(c *fiber.Ctx) error
	GetWarehouseProductByWarehouseIDAndProductID(c *fiber.Ctx) error
	UpdateWarehouseProduct(c *fiber.Ctx) error
	DeleteWarehouseProduct(c *fiber.Ctx) error
	DeleteAllWarehouseProductByProductID(c *fiber.Ctx) error
	GetWarehouseProductByProductID(c *fiber.Ctx) error
	GetProductTotalStock(c *fiber.Ctx) error
	GetDetailWarehouse(c *fiber.Ctx) error
	GetDetailWarehouseProductByID(c *fiber.Ctx) error
}

type WarehouseProductController struct {
	warehouseProductUsecase usecase.WarehouseProductUsecaseInterface
}

func NewWarehouseProductController(warehouseProductUsecase usecase.WarehouseProductUsecaseInterface) WarehouseProductControllerInterface {
	return &WarehouseProductController{
		warehouseProductUsecase: warehouseProductUsecase,
	}
}

// CreateWarehouseProduct implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) CreateWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.CreateWarehouseProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouseProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouseProduct - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	warehouseID := conv.StringToUint(c.Params("warehouse_id"))

	warehouseProductModel := model.WarehouseProduct{
		WarehouseID: warehouseID,
		ProductID:   req.ProductID,
		Stock:       req.Stock,
	}

	if err := w.warehouseProductUsecase.CreateWarehouseProduct(ctx, &warehouseProductModel); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouseProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Warehouse product created successfully",
	})
}

// DeleteAllWarehouseProductByProductID implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) DeleteAllWarehouseProductByProductID(c *fiber.Ctx) error {
	ctx := c.Context()

	productID := conv.StringToUint(c.Params("product_id"))

	if err := w.warehouseProductUsecase.DeleteAllWarehouseProductByProductID(ctx, productID); err != nil {
		log.Errorf("[WarehouseProductController] DeleteAllWarehouseProductByProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All warehouse products for the given product ID deleted successfully",
	})
}

// DeleteWarehouseProduct implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) DeleteWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	WarehouseProductID := conv.StringToUint(c.Params("warehouse_product_id"))

	if err := w.warehouseProductUsecase.DeleteWarehouseProduct(ctx, WarehouseProductID); err != nil {
		log.Errorf("[WarehouseProductController] DeleteWarehouseProduct - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse product deleted successfully",
	})
}

// GetDetailWarehouse implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) GetDetailWarehouse(c *fiber.Ctx) error {
	ctx := c.Context()

	warehouseID := conv.StringToUint(c.Params("warehouse_id"))

	warehouse, products, err := w.warehouseProductUsecase.GetDetailWarehouse(ctx, warehouseID)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetDetailWarehouse - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := response.DetailWarehouseResponse{
		ID:      warehouse.ID,
		Name:    warehouse.Name,
		Address: warehouse.Address,
		Photo:   warehouse.Photo,
		Phone:   warehouse.Phone,
	}

	productMap := make(map[uint]*httpclient.ProductResponse)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	for _, wp := range warehouse.WarehouseProducts {
		warehouseProduct := response.WarehouseProductResponse{
			ID:          wp.ID,
			ProductID:   wp.ProductID,
			WarehouseID: wp.WarehouseID,
			Stock:       wp.Stock,
		}

		if product, exist := productMap[wp.ProductID]; exist {
			warehouseProduct.ProductName = product.Name
			warehouseProduct.ProductAbout = product.About
			warehouseProduct.ProductPhoto = product.Thumbnail
			warehouseProduct.ProductPrice = int(product.Price)
			warehouseProduct.ProductCategory = product.Category.Name
			warehouseProduct.ProductCategoryPhoto = product.Category.Photo
		}

		resp.WarehouseProducts = append(resp.WarehouseProducts, warehouseProduct)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Warehouse detail retrieved successfully",
	})
}

// GetDetailWarehouseProductByID implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) GetDetailWarehouseProductByID(c *fiber.Ctx) error {
	ctx := c.Context()

	warehouseProductID := conv.StringToUint(c.Params("warehouse_product_id"))

	warehouseProduct, product, err := w.warehouseProductUsecase.GetDetailWarehouseProductByID(ctx, warehouseProductID)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetDetailWarehouseProductByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := response.GetDetailWarehouseProductResponse{
		ID:               warehouseProduct.ID,
		WarehouseID:      warehouseProduct.WarehouseID,
		ProductID:        warehouseProduct.ProductID,
		Stock:            warehouseProduct.Stock,
		WarehouseName:    warehouseProduct.Warehouse.Name,
		WarehousePhoto:   warehouseProduct.Warehouse.Photo,
		WarehousePhone:   warehouseProduct.Warehouse.Phone,
		ProductName:      product.Name,
		ProductBarcode:   product.Barcode,
		ProductPrice:     int(product.Price),
		ProductAbout:     product.About,
		ProductThumbnail: product.Thumbnail,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Warehouse product detail retrieved successfully",
	})
}

// GetProductTotalStock implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) GetProductTotalStock(c *fiber.Ctx) error {
	ctx := c.Context()

	productID := conv.StringToUint(c.Params("product_id"))

	totalStock, err := w.warehouseProductUsecase.GetProductTotalStock(ctx, productID)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetProductTotalStock - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := response.ProductTotalStockResponse{
		ProductID:  productID,
		TotalStock: totalStock,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Product total stock fetched successfully",
	})
}

// GetWarehouseProductByProductID implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) GetWarehouseProductByProductID(c *fiber.Ctx) error {
	ctx := c.Context()

	productID := conv.StringToUint(c.Params("product_id"))

	warehouseProduct, err := w.warehouseProductUsecase.GetWarehouseProductByProductID(ctx, productID)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetWarehouseProductByProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.WarehouseResponse{}
	for _, wp := range warehouseProduct {
		resp = append(resp, response.WarehouseResponse{
			ID:      wp.Warehouse.ID,
			Name:    wp.Warehouse.Name,
			Address: wp.Warehouse.Address,
			Photo:   wp.Warehouse.Photo,
			Phone:   wp.Warehouse.Phone,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Warehouse products retrieved successfully",
	})
}

// GetWarehouseProductByWarehouseIDAndProductID implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) GetWarehouseProductByWarehouseIDAndProductID(c *fiber.Ctx) error {
	ctx := c.Context()

	warehouseID := conv.StringToUint(c.Params("warehouse_id"))
	productID := conv.StringToUint(c.Params("product_id"))

	warehouseProduct, err := w.warehouseProductUsecase.GetWarehouseProductByWarehouseIDAndProductID(ctx, warehouseID, productID)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetWarehouseProductByWarehouseIDAndProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := response.WarehouseProductResponse{
		ID:          warehouseProduct.ID,
		WarehouseID: warehouseProduct.WarehouseID,
		ProductID:   warehouseProduct.ProductID,
		Stock:       warehouseProduct.Stock,
		Warehouse: response.WarehouseResponse{
			ID:      warehouseProduct.Warehouse.ID,
			Name:    warehouseProduct.Warehouse.Name,
			Address: warehouseProduct.Warehouse.Address,
			Photo:   warehouseProduct.Warehouse.Photo,
			Phone:   warehouseProduct.Warehouse.Phone,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Warehouse product retrieved successfully",
	})
}

// UpdateWarehouseProduct implements [WarehouseProductControllerInterface].
func (w *WarehouseProductController) UpdateWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	warehouseProductID := conv.StringToUint(c.Params("warehouse_product_id"))
	warehouseID := conv.StringToUint(c.Params("warehouse_id"))

	req := request.CreateWarehouseProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.WarehouseProduct{
		ID:          warehouseProductID,
		WarehouseID: warehouseID,
		ProductID:   req.ProductID,
		Stock:       req.Stock,
	}

	if err := w.warehouseProductUsecase.UpdateWarehouseProduct(ctx, &reqModel); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update warehouse product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse product updated successfully",
	})
}
