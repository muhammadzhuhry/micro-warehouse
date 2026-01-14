package model

import "time"

type Merchant struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null"`
	Address   string     `json:"address"`
	Photo     string     `json:"photo"`
	Phone     string     `json:"phone"`
	KeeperID  uint64     `json:"keeper_id" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	MerchantProducts []MerchantProduct `json:"merchant_products,omitempty" gorm:"foreignKey:MerchantID;references:ID"`
}
