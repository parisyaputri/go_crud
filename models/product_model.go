package models

import "time"

type Product struct {
	ID          int
	ProductName string
	Price       int
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
