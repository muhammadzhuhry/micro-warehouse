package app

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *Container) {
	api := app.Group("/api/v1")

	warehouses := api.Group("/warehouses")
	warehouses.Post("/", container.WarehouseController.CreateWarehouse)
	warehouses.Get("/", container.WarehouseController.GetAllWarehouses)
	warehouses.Get("/:id", container.WarehouseController.GetWarehouseByID)
	warehouses.Put("/:id", container.WarehouseController.UpdateWarehouse)
	warehouses.Delete("/:id", container.WarehouseController.DeleteWarehouse)

	warehouseProduct := api.Group("/warehouse-products")
	warehouseProduct.Post("/:warehouse_id", container.WarehouseProductController.CreateWarehouseProduct)
	warehouseProduct.Get("/:warehouse_id", container.WarehouseProductController.GetDetailWarehouse)
	warehouseProduct.Get("/:warehouse_id/detail/:product_id", container.WarehouseProductController.GetWarehouseProductByWarehouseIDAndProductID)
	warehouseProduct.Put("/detail/:warehouse_product_id", container.WarehouseProductController.UpdateWarehouseProduct)
	warehouseProduct.Delete("/detail/:warehouse_product_id", container.WarehouseProductController.DeleteWarehouseProduct)
	warehouseProduct.Delete("/detail/products/:product_id", container.WarehouseProductController.DeleteAllWarehouseProductByProductID)
	warehouseProduct.Get("/detail/products/:product_id/total-stock", container.WarehouseProductController.GetProductTotalStock)
	warehouseProduct.Get("/detail/products/:product_id", container.WarehouseProductController.GetWarehouseProductByProductID)
	warehouseProduct.Get("/detail/products/:product_id/warehouses", container.WarehouseProductController.GetDetailWarehouseProductByID)

	api.Post("/uploads-warehouse", container.UploadController.UplodaPhoto)
}
