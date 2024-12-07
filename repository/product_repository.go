package repository

import (
	"database/sql"
	"go-crud/models"
	"time"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) (*models.Product, error)
	FindProductByID(ID int) (*models.Product, error)
	FindAllProduct() ([]*models.Product, error)
	UpdateProduct(ID int, product *models.Product) (*models.Product, error)
	DeleteProduct(ID int) error // Tambahkan ini

}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{db}
}

// CreateProduct inserts a new product into the database
func (r *productRepository) CreateProduct(product *models.Product) (*models.Product, error) {
	query := "INSERT INTO products (ProductName, Price, Quantity, CreatedAt, UpdatedAt) VALUES (?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, product.ProductName, product.Price, product.Quantity, time.Now(), time.Now())

	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	product.ID = int(lastID)
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	return product, nil
}

// FindProductByID retrieves a product by its ID
func (r *productRepository) FindProductByID(ID int) (*models.Product, error) {
	query := "SELECT id, ProductName, Price, Quantity, CreatedAt, UpdatedAt FROM products WHERE id = ?"
	row := r.db.QueryRow(query, ID)

	product := &models.Product{}
	err := row.Scan(&product.ID, &product.ProductName, &product.Price, &product.Quantity, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return product, nil
}

// FindAllProduct retrieves all products from the database
func (r *productRepository) FindAllProduct() ([]*models.Product, error) {
	query := "SELECT id, ProductName, Price, Quantity, CreatedAt, UpdatedAt FROM products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(&product.ID, &product.ProductName, &product.Price, &product.Quantity, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct updates an existing product in the database
func (r *productRepository) UpdateProduct(ID int, product *models.Product) (*models.Product, error) {
	query := "UPDATE products SET ProductName = ?, Price = ?, Quantity = ?, UpdatedAt = ? WHERE id = ?"
	_, err := r.db.Exec(query, product.ProductName, product.Price, product.Quantity, time.Now(), ID)
	if err != nil {
		return nil, err
	}

	product.ID = ID
	product.UpdatedAt = time.Now()

	return product, nil
}

// DeleteProduct deletes a product by its ID
func (r *productRepository) DeleteProduct(ID int) error {
	query := "DELETE FROM products WHERE id = ?"
	result, err := r.db.Exec(query, ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
