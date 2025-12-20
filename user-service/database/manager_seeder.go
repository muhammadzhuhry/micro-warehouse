package database

import (
	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/pkg/conv"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func SeedManager(db *gorm.DB) {
	bytes, err := conv.HashPassword("manager123")
	if err != nil {
		log.Fatalf("%s : %v", err.Error(), err)
	}
	
	modelRole := model.Role{}

	// Check if Manager role exists
	if err := db.Where("name = ?", "Manager").First(&modelRole).Error; err != nil {
		log.Fatalf("%s : %v", err.Error(), err)
	}

	// Mapping admin user
	admin := model.User{
		Name: "manager",
		Email: "manager@email.com",
		Password: bytes,
		Roles: []model.Role{modelRole},
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "manager@email.com"}).Error; err != nil {
		log.Fatalf("%s : %v", err.Error(), err)
	} else {
		log.Infof("Admin %s user created", admin.Name)
	}
}