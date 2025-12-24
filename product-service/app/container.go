package app

import (
	"micro-warehouse/product-service/configs"
	"micro-warehouse/product-service/controller"
	"micro-warehouse/product-service/database"
	"micro-warehouse/product-service/pkg/storage"
	"micro-warehouse/product-service/repository"
	"micro-warehouse/product-service/usecase"

	"github.com/gofiber/fiber/v2/log"
)

type Container struct {
	CategoryController controller.CategoryControllerInterface
	ProductController  controller.ProductControllerInterface
	UploadController   controller.UploadControllerInterface
}

func BuildContainer() *Container {
	config := configs.NewConfig()

	// Initialize database connection
	db, err := database.ConnectionPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Initialize storage service and upload helper
	supabaseStorage := storage.NewSupabaseStorage(*config)
	fileUploadHelper := storage.NewFileUploadHelper(supabaseStorage, *config)

	// Initialize repositories, usecases, and controllers

	categoryRepo := repository.NewCategoryRepository(db.DB)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryController := controller.NewCategoryController(categoryUsecase)

	productRepo := repository.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productController := controller.NewProductController(productUsecase)

	uploadController := controller.NewUploadController(fileUploadHelper)

	return &Container{
		CategoryController: categoryController,
		ProductController:  productController,
		UploadController:   uploadController,
	}
}
