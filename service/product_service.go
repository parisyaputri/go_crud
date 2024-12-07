package service

import (
	"go-crud/input"
	"go-crud/models"
	"go-crud/repository"
)

type ProductService interface {
	CreateProduct(input *input.ProductInput) (*models.Product, error)
	GetProductByID(ID int) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	UpdateProduct(ID int, input *input.ProductInput) (*models.Product, error)
	DeleteProduct(ID int) error
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo}
}

func (s *productService) CreateProduct(input *input.ProductInput) (*models.Product, error) {
	product := &models.Product{
		ProductName: input.ProductName,
		Price:       input.Price,
		Quantity:    input.Quantity,
	}

	return s.productRepo.CreateProduct(product)
}

func (s *productService) GetProductByID(ID int) (*models.Product, error) {
	return s.productRepo.FindProductByID(ID)
}

func (s *productService) GetAllProducts() ([]*models.Product, error) {
	return s.productRepo.FindAllProduct()
}

func (s *productService) UpdateProduct(ID int, input *input.ProductInput) (*models.Product, error) {
	product := &models.Product{
		ProductName: input.ProductName,
		Price:       input.Price,
		Quantity:    input.Quantity,
	}

	return s.productRepo.UpdateProduct(ID, product)
}

func (s *productService) DeleteProduct(ID int) error {
	return s.productRepo.DeleteProduct(ID)
}
