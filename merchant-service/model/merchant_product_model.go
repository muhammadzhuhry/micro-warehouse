package model

import "time"

type MerchantProduct struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID  int64      `json:"merchant_id" gorm:"not null;index"`
	ProductID   int64      `json:"product_id" gorm:"not null;index"`
	WarehouseID int64      `json:"warehouse_id" gorm:"not null;index"`
	Stock       int        `json:"stock" gorm:"not null;default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`

	Merchant Merchant `json:"merchant,omitempty" gorm:"foreignKey:MerchantID;references:ID"`
}
