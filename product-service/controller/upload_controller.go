package controller

import (
	"micro-warehouse/product-service/controller/response"
	"micro-warehouse/product-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadProductImage(c *fiber.Ctx) error
	UploadCategoryImage(c *fiber.Ctx) error
}

type UploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

func NewUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &UploadController{
		fileUploadHelper: fileUploadHelper,
	}
}

// UploadCategoryImage implements [UploadControllerInterface].
func (u *UploadController) UploadCategoryImage(c *fiber.Ctx) error {
	ctx := c.Context()

	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("[UploadController] UploadCategoryImage - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to get image from form data",
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(ctx, file)
	if err != nil {
		log.Errorf("[UploadController] UploadCategoryImage - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to upload image",
		})
	}

	resp := response.UploadResponse{
		URL:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "image uploaded successfully",
	})
}

// UploadProductImage implements [UploadControllerInterface].
func (u *UploadController) UploadProductImage(c *fiber.Ctx) error {
	ctx := c.Context()

	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("[UploadController] UploadProductImage - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to get image from form data",
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(ctx, file)
	if err != nil {
		log.Errorf("[UploadController] UploadProductImage - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to upload image",
		})
	}

	resp := response.UploadResponse{
		URL:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "image uploaded successfully",
	})
}
