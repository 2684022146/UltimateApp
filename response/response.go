package response

import "github.com/shopspring/decimal"

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type ProductResponse struct {
	Name          string          `json:"name"`
	Price         decimal.Decimal `json:"price"`
	Description   string          `json:"description"`
	Sku           string          `json:"sku"`
	ProductJSONLD ProductJSONLD   `json:"product_ld"`
}
type ProductJSONLD struct {
	Context string          `json:"@context"`
	Type    string          `json:"@type"`
	Name    string          `json:"name"`
	Price   decimal.Decimal `json:"price"`
}
