package controller

import (
	"micro-warehouse/product-service/controller/request"
	"micro-warehouse/product-service/controller/response"
	"micro-warehouse/product-service/model"
	"micro-warehouse/product-service/pkg/conv"
	"micro-warehouse/product-service/pkg/pagination"
	"micro-warehouse/product-service/usecase"
	"micro-warehouse/product-service/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ProductControllerInterface interface {
	CreateProduct(c *fiber.Ctx) error
	GetAllProducts(c *fiber.Ctx) error
	GetProductByID(c *fiber.Ctx) error
	GetProductByBarcode(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type ProductController struct {
	productUsecase usecase.ProductUsecaseInterface
}

func NewProductController(productUsecase usecase.ProductUsecaseInterface) ProductControllerInterface {
	return &ProductController{productUsecase: productUsecase}
}

// CreateProduct implements [ProductControllerInterface].
func (p *ProductController) CreateProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.CreateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[ProductController] CreateProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[ProductController] CreateProduct - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	productModel := model.Product{
		Name:       req.Name,
		Barcode:    req.Barcode,
		CategoryID: req.CategoryID,
		Thumbnail:  req.Thumbnail,
		About:      req.About,
		Price:      conv.IntToFloat64(req.Price),
		IsPopular:  req.IsPopular,
	}

	if err := p.productUsecase.CreateProduct(ctx, productModel); err != nil {
		log.Errorf("[ProductController] CreateProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
	})
}

// DeleteProduct implements [ProductControllerInterface].
func (p *ProductController) DeleteProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	productID := c.Params("id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product id is required",
		})
	}

	id := conv.StringToUint(productID)

	if err := p.productUsecase.DeleteProduct(ctx, id); err != nil {
		log.Errorf("[ProductController] DeleteProduct - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

// GetAllProducts implements [ProductControllerInterface].
func (p *ProductController) GetAllProducts(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.GetAllProductsRequest{}
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[ProductController] GetAllProducts - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[ProductController] GetAllProducts - 2: %v", err)
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

	products, totalRecords, err := p.productUsecase.GetAllProducts(ctx, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[ProductController] GetAllProducts - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get products",
		})
	}

	resp := []response.ProductResponse{}
	for _, product := range products {
		resp = append(resp, response.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			Barcode:    product.Barcode,
			CategoryID: product.Category.ID,
			Category: response.CategoryResponse{
				ID:            product.Category.ID,
				Name:          product.Category.Name,
				Tagline:       product.Category.Tagline,
				Photo:         product.Category.Photo,
				CountProducts: len(product.Category.Products),
			},
			Thumbnail: product.Thumbnail,
			About:     product.About,
			Price:     product.Price,
			IsPopular: product.IsPopular,
		})
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(totalRecords))

	response := response.GetAllProductsResponse{
		Products:   resp,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Products retrieved successfully",
	})
}

// GetProductByBarcode implements [ProductControllerInterface].
func (p *ProductController) GetProductByBarcode(c *fiber.Ctx) error {
	ctx := c.Context()

	barcode := c.Params("barcode")
	if barcode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product barcode is required",
		})
	}

	product, err := p.productUsecase.GetProductByBarcode(ctx, barcode)
	if err != nil {
		log.Errorf("[ProductController] GetProductByBarcode - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get product",
		})
	}

	resp := response.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Barcode:    product.Barcode,
		About:      product.About,
		Price:      product.Price,
		IsPopular:  product.IsPopular,
		CategoryID: product.Category.ID,
		Thumbnail:  product.Thumbnail,
		Category: response.CategoryResponse{
			ID:            product.Category.ID,
			Name:          product.Category.Name,
			Tagline:       product.Category.Tagline,
			Photo:         product.Category.Photo,
			CountProducts: len(product.Category.Products),
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Product retrieved successfully",
	})
}

// GetProductByID implements [ProductControllerInterface].
func (p *ProductController) GetProductByID(c *fiber.Ctx) error {
	ctx := c.Context()

	productID := c.Params("id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product id is required",
		})
	}

	id := conv.StringToUint(productID)

	product, err := p.productUsecase.GetProductByID(ctx, id)
	if err != nil {
		log.Errorf("[ProductController] GetProductByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get product",
		})
	}

	resp := response.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Barcode:    product.Barcode,
		About:      product.About,
		Price:      product.Price,
		IsPopular:  product.IsPopular,
		CategoryID: product.Category.ID,
		Thumbnail:  product.Thumbnail,
		Category: response.CategoryResponse{
			ID:            product.Category.ID,
			Name:          product.Category.Name,
			Tagline:       product.Category.Tagline,
			Photo:         product.Category.Photo,
			CountProducts: len(product.Category.Products),
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Product retrieved successfully",
	})
}

// UpdateProduct implements [ProductControllerInterface].
func (p *ProductController) UpdateProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.CreateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[ProductController] UpdateProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[ProductController] UpdateProduct - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	productModel := model.Product{
		ID:         conv.StringToUint(c.Params("id")),
		Name:       req.Name,
		Barcode:    req.Barcode,
		CategoryID: req.CategoryID,
		Thumbnail:  req.Thumbnail,
		About:      req.About,
		Price:      conv.IntToFloat64(req.Price),
		IsPopular:  req.IsPopular,
	}

	if err := p.productUsecase.UpdateProduct(ctx, productModel); err != nil {
		log.Errorf("[ProductController] UpdateProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}
