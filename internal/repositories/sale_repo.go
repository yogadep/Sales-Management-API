package repositories

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sales-management-api/internal/models"
)

type SaleRepo struct {
	db *gorm.DB
}

func NewSaleRepo(db *gorm.DB) *SaleRepo {
	return &SaleRepo{db: db}
}

type SaleItemInput struct {
	ProductID uint
	Qty       int64
}

func (r *SaleRepo) CreateSale(cashierID uint, items []SaleItemInput) (*models.Sale, error) {
	if len(items) == 0 {
		return nil, errors.New("items cannot be empty")
	}

	var createdSale models.Sale

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// create sale header (total sementara 0)
		sale := models.Sale{
			CashierID: cashierID,
			Total:     0,
		}
		if err := tx.Create(&sale).Error; err != nil {
			return err
		}

		var total int64 = 0

		// process items
		for _, it := range items {
			if it.Qty <= 0 {
				return errors.New("qty must be > 0")
			}

			// Lock row produk biar aman dari race condition (stok minus karena transaksi paralel)
			var p models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&p, it.ProductID).Error; err != nil {
				return err
			}

			if p.Stock < it.Qty {
				return errors.New("insufficient stock")
			}

			// reduce stock
			p.Stock = p.Stock - it.Qty
			if err := tx.Save(&p).Error; err != nil {
				return err
			}

			subtotal := p.Price * it.Qty
			total += subtotal

			saleItem := models.SaleItem{
				SaleID:    sale.ID,
				ProductID: p.ID,
				Qty:       it.Qty,
				Price:     p.Price, // snapshot
				Subtotal:  subtotal,
			}
			if err := tx.Create(&saleItem).Error; err != nil {
				return err
			}
		}

		// update sale total
		if err := tx.Model(&models.Sale{}).
			Where("id = ?", sale.ID).
			Update("total", total).Error; err != nil {
			return err
		}

		// preload items buat response
		if err := tx.Preload("Items").First(&createdSale, sale.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &createdSale, nil
}

func (r *SaleRepo) FindAll(limit int) ([]models.Sale, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var sales []models.Sale
	err := r.db.
		Preload("Items").
		Order("id desc").
		Limit(limit).
		Find(&sales).Error
	return sales, err
}

func (r *SaleRepo) FindByID(id uint) (*models.Sale, error) {
	var sale models.Sale
	if err := r.db.
		Preload("Items").
		First(&sale, id).Error; err != nil {
		return nil, err
	}
	return &sale, nil
}
