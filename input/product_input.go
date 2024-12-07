package input

type ProductInput struct {
	ProductName string `json:"product_name" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
}
