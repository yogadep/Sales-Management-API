package services

import (
	"sales-management-api/internal/models"
	"sales-management-api/internal/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepo
}

func NewProductService(r *repositories.ProductRepo) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) Create(p *models.Product) error {
	return s.repo.Create(p)
}

func (s *ProductService) Update(id uint, p *models.Product) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	existing.Name = p.Name
	existing.SKU = p.SKU
	existing.Price = p.Price
	existing.Stock = p.Stock

	return s.repo.Update(existing)
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ProductService) List() ([]models.Product, error) {
	return s.repo.FindAll()
}
