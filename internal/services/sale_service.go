package services

import (
	"sales-management-api/internal/models"
	"sales-management-api/internal/repositories"
)

type SaleService struct {
	repo *repositories.SaleRepo
}

func NewSaleService(r *repositories.SaleRepo) *SaleService {
	return &SaleService{repo: r}
}

func (s *SaleService) Create(cashierID uint, items []repositories.SaleItemInput) (*models.Sale, error) {
	return s.repo.CreateSale(cashierID, items)
}

func (s *SaleService) List(limit int) ([]models.Sale, error) {
	return s.repo.FindAll(limit)
}

func (s *SaleService) Detail(id uint) (*models.Sale, error) {
	return s.repo.FindByID(id)
}
