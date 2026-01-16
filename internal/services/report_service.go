package services

import (
	"time"

	"sales-management-api/internal/models"
	"sales-management-api/internal/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepo
}

func NewReportService(r *repositories.ReportRepo) *ReportService {
	return &ReportService{repo: r}
}

func (s *ReportService) Sales(from, to time.Time) ([]models.Sale, error) {
	return s.repo.SalesReport(from, to)
}
