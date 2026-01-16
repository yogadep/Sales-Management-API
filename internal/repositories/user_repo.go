package repositories

import (
	"sales-management-api/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// =====================
// AUTH
// =====================
func (r *UserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// =====================
// ADMIN / USER MANAGEMENT
// =====================

// List all users (WITHOUT password)
func (r *UserRepo) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.
		Select("id", "username", "role", "created_at").
		Order("id asc").
		Find(&users).Error
	return users, err
}

// Detail user by ID (WITHOUT password)
func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.
		Select("id", "username", "role", "created_at").
		First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
