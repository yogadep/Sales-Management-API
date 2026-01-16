package repositories

import (
	"sales-management-api/internal/models"

	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(p *models.Product) error {
	return r.db.Create(p).Error
}

func (r *ProductRepo) Update(p *models.Product) error {
	return r.db.Save(p).Error
}

func (r *ProductRepo) Delete(id uint) error {
	res := r.db.Delete(&models.Product{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ProductRepo) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Order("id desc").Find(&products).Error
	return products, err
}

func (r *ProductRepo) FindByID(id uint) (*models.Product, error) {
	var p models.Product
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
