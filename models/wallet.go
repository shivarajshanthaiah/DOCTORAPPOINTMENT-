package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User Patient
	Amount float64 `json:"amount"`
}
