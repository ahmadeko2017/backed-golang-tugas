package repository

import (
	"github.com/ahmadeko2017/backed-golang-tugas/internal/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	FindAll(name string, page int, pageSize int) ([]entity.Product, int64, error)
	FindByID(id uint) (entity.Product, error)
	FindByIDWithLock(tx *gorm.DB, id uint) (entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll(name string, page int, pageSize int) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.Model(&entity.Product{}).Preload("Category")

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&products).Error

	return products, total, err
}

func (r *productRepository) FindByID(id uint) (entity.Product, error) {
	var product entity.Product
	err := r.db.Preload("Category").First(&product, id).Error
	return product, err
}

func (r *productRepository) FindByIDWithLock(tx *gorm.DB, id uint) (entity.Product, error) {
	var product entity.Product
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, id).Error; err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepository) Update(product *entity.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Product{}, id).Error
}
