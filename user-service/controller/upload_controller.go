package controller

import (
	"micro-warehouse/user-service/controller/response"
	"micro-warehouse/user-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadPhoto(c *fiber.Ctx) error
}

type UploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

func NewUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &UploadController{
		fileUploadHelper: fileUploadHelper,
	}
}

// UploadPhoto implements [UploadControllerInterface].
func (u *UploadController) UploadPhoto(c *fiber.Ctx) error {
	ctx := c.Context()

	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("[UploadController] UploadPhoto - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to get image from form data",
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(ctx, file)
	if err != nil {
		log.Errorf("[UploadController] UploadPhoto - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to upload image",
		})
	}

	resp := response.UploadPhotoResponse{
		URL:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "image uploaded successfully",
	})
}
