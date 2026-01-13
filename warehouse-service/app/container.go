package app

import (
	"log"
	"micro-warehouse/warehouse-service/configs"
	"micro-warehouse/warehouse-service/controller"
	"micro-warehouse/warehouse-service/database"
	"micro-warehouse/warehouse-service/pkg/httpclient"
	"micro-warehouse/warehouse-service/pkg/storage"
	"micro-warehouse/warehouse-service/repository"
	"micro-warehouse/warehouse-service/usecase"
)

type Container struct {
	WarehouseController        controller.WarehouseControllerInterface
	WarehouseProductController controller.WarehouseProductControllerInterface
	UploadController           controller.UploadControllerInterface
}

func BuildContainer() *Container {
	config := configs.NewConfig()

	// Initialize database connection
	db, err := database.ConnectionPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	supabaseStorage := storage.NewSupabaseStorage(*config)
	fileUploadHelper := storage.NewFileUploadHelper(supabaseStorage, *config)

	warehouseRepo := repository.NewWarehouseRepository(db.DB)
	warehouseUsecase := usecase.NewWarehouseUsecase(warehouseRepo)
	warehouseController := controller.NewWarehouseController(warehouseUsecase)

	warehouseProductRepo := repository.NewWarehouseProductRepository(db.DB)
	productClient := httpclient.NewProductClient(*config)
	warehouseProductUsecase := usecase.NewWarehouseProductUsecase(warehouseProductRepo, productClient)
	warehouseProductController := controller.NewWarehouseProductController(warehouseProductUsecase)

	uploadController := controller.NewUploadController(fileUploadHelper)

	return &Container{
		WarehouseController:        warehouseController,
		WarehouseProductController: warehouseProductController,
		UploadController:           uploadController,
	}
}
