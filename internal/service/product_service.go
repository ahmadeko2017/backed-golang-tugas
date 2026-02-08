package service

import (
	"errors"

	"github.com/ahmadeko2017/backed-golang-tugas/internal/entity"
	"github.com/ahmadeko2017/backed-golang-tugas/internal/repository"
)

type ProductService interface {
	CreateProduct(product *entity.Product) error
	GetAllProducts(name string, page int, limit int) ([]entity.Product, int64, error)
	GetProductByID(id uint) (entity.Product, error)
	UpdateProduct(id uint, product *entity.Product) error
	DeleteProduct(id uint) error
}

type productService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(repo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) CreateProduct(product *entity.Product) error {
	// Validate Category ID
	_, err := s.categoryRepo.FindByID(product.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}
	return s.repo.Create(product)
}

func (s *productService) GetAllProducts(name string, page int, limit int) ([]entity.Product, int64, error) {
	return s.repo.FindAll(name, page, limit)
}

func (s *productService) GetProductByID(id uint) (entity.Product, error) {
	return s.repo.FindByID(id)
}

func (s *productService) UpdateProduct(id uint, product *entity.Product) error {
	existingProduct, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Update fields
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.Stock = product.Stock

	// Optional: Validate Category ID if changed
	if product.CategoryID != 0 {
		_, err := s.categoryRepo.FindByID(product.CategoryID)
		if err != nil {
			return errors.New("category not found")
		}
		existingProduct.CategoryID = product.CategoryID
	}

	return s.repo.Update(&existingProduct)
}

func (s *productService) DeleteProduct(id uint) error {
	return s.repo.Delete(id)
}
