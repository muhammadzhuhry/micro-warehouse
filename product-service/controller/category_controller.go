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

type CategoryControllerInterface interface {
	CreateCategory(ctx *fiber.Ctx) error
	GetAllCategories(ctx *fiber.Ctx) error
	GetCategoryByID(ctx *fiber.Ctx) error
	UpdateCategory(ctx *fiber.Ctx) error
	DeleteCategory(ctx *fiber.Ctx) error
}

type categoryController struct {
	categoryUsecase usecase.CategoryUsecaseInterface
}

func NewCategoryController(categoryUsecase usecase.CategoryUsecaseInterface) CategoryControllerInterface {
	return &categoryController{categoryUsecase: categoryUsecase}
}

// CreateCategory implements [CategoryControllerInterface].
func (c *categoryController) CreateCategory(ctx *fiber.Ctx) error {
	cx := ctx.Context()

	req := request.CreateCategoryRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[CategoryController] CreateCategory - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[CategoryController] CreateCategory - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	categoryModel := model.Category{
		Name:    req.Name,
		Tagline: req.Tagline,
		Photo:   req.Photo,
	}

	if err := c.categoryUsecase.CreateCategory(cx, categoryModel); err != nil {
		log.Errorf("[CategoryController] CreateCategory - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create category",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
	})
}

// DeleteCategory implements [CategoryControllerInterface].
func (c *categoryController) DeleteCategory(ctx *fiber.Ctx) error {
	cx := ctx.Context()

	categoryID := ctx.Params("id")
	if categoryID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Category id is required",
		})
	}

	id := conv.StringToUint(categoryID)

	if err := c.categoryUsecase.DeleteCategory(cx, id); err != nil {
		log.Errorf("[CategoryController] DeleteCategory - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete category",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}

// GetAllCategories implements [CategoryControllerInterface].
func (c *categoryController) GetAllCategories(ctx *fiber.Ctx) error {
	cx := ctx.Context()

	req := request.GetAllCategoriesRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		log.Errorf("[CategoryController] GetAllCategories - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[CategoryController] GetAllCategories - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	categories, totalRecords, err := c.categoryUsecase.GetAllCategories(cx, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[CategoryController] GetAllCategories - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.CategoryResponse{}
	for _, category := range categories {
		resp = append(resp, response.CategoryResponse{
			ID:      category.ID,
			Name:    category.Name,
			Tagline: category.Tagline,
			Photo:   category.Photo,
		})
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(totalRecords))

	response := response.GetAllCategoriesResponse{
		Categories: resp,
		Pagination: paginationInfo,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Categories retrieved successfully",
	})
}

// GetCategoryByID implements [CategoryControllerInterface].
func (c *categoryController) GetCategoryByID(ctx *fiber.Ctx) error {
	cx := ctx.Context()

	categoryID := ctx.Params("id")
	if categoryID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Category id is required",
		})
	}

	id := conv.StringToUint(categoryID)

	category, err := c.categoryUsecase.GetCategoryByID(cx, id)
	if err != nil {
		log.Errorf("[CategoryController] GetCategoryByID - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get category",
		})
	}

	resp := response.CategoryResponse{
		ID:            category.ID,
		Name:          category.Name,
		Tagline:       category.Tagline,
		Photo:         category.Photo,
		CountProducts: len(category.Products),
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Category retrieved successfully",
	})
}

// UpdateCategory implements [CategoryControllerInterface].
func (c *categoryController) UpdateCategory(ctx *fiber.Ctx) error {
	cx := ctx.Context()

	req := request.CreateCategoryRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[CategoryController] UpdateCategory - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[CategoryController] UpdateCategory - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	categoryModel := model.Category{
		ID:      conv.StringToUint(ctx.Params("id")),
		Name:    req.Name,
		Tagline: req.Tagline,
		Photo:   req.Photo,
	}

	if err := c.categoryUsecase.UpdateCategory(cx, categoryModel); err != nil {
		log.Errorf("[CategoryController] UpdateCategory - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update category",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category updated successfully",
	})
}
