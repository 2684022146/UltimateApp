package repository

import (
	"demo01/model"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductList() ([]*model.Product, error)
	GetProductDetail(name string) (*model.Product, error)
}
type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) GetProductList() ([]*model.Product, error) {
	var productList []*model.Product
	if err := r.db.Table("my_product").Find(&productList).Error; err != nil {
		return nil, fmt.Errorf("query fail:%s", err)
	}
	if len(productList) == 0 {
		log.Println("no data")
		return []*model.Product{}, nil
	}

	return productList, nil
}
func (r *productRepository) GetProductDetail(sku string) (*model.Product, error) {
	if sku == "" {
		return nil, nil
	}
	var product *model.Product
	if err := r.db.Table("my_product").Where("sku=?", sku).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未找到:%s", err)
		}
		return nil, fmt.Errorf("query fail:%s", err)
	}
	fmt.Println(product.Name, product.Price)
	return product, nil
}
