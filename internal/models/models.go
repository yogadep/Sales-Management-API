package models

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleKasir Role = "kasir"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;size:80;not null"`
	PasswordHash string `gorm:"not null"`
	Role         Role   `gorm:"size:10;not null"`
	CreatedAt    time.Time
}

type Product struct {
	ID        uint   `gorm:"primaryKey"`
	SKU       string `gorm:"uniqueIndex;size:50;not null"`
	Name      string `gorm:"size:150;not null"`
	Price     int64  `gorm:"not null"`
	Stock     int64  `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Sale struct {
	ID        uint  `gorm:"primaryKey"`
	CashierID uint  `gorm:"not null"`
	Total     int64 `gorm:"not null"`
	CreatedAt time.Time
}

type SaleItem struct {
	ID        uint  `gorm:"primaryKey"`
	SaleID    uint  `gorm:"index;not null"`
	ProductID uint  `gorm:"index;not null"`
	Qty       int64 `gorm:"not null"`
	Price     int64 `gorm:"not null"`
	Subtotal  int64 `gorm:"not null"`
}
