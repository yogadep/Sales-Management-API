package repositories

import (
	"time"

	"gorm.io/gorm"

	"sales-management-api/internal/models"
)

type ReportRepo struct {
	db *gorm.DB
}

func NewReportRepo(db *gorm.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

func (r *ReportRepo) SalesReport(from, to time.Time) ([]models.Sale, error) {
	var sales []models.Sale
	err := r.db.
		Preload("Items").
		Where("created_at BETWEEN ? AND ?", from, to).
		Order("created_at asc").
		Find(&sales).Error
	return sales, err
}
