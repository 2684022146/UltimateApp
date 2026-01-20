package model

import "github.com/shopspring/decimal"

type Product struct {
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	Sku         string          `json:"sku"`
}
