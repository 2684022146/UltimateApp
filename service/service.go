package service

import (
	"context"
	"webdemo/model"
	"webdemo/repository"

	"fmt"
)

type ProductService interface {
	ProductList(ctx context.Context) ([]*model.Product, error)
	ProductDetail(ctx context.Context, name string) (*model.Product, error)
}
type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}
func (s *productService) ProductList(ctx context.Context) ([]*model.Product, error) {
	productList, err := s.repo.GetProductList()
	if err != nil {
		return nil, fmt.Errorf("query all product fail:%w", err)
	}
	return productList, nil
}
func (s *productService) ProductDetail(ctx context.Context, sku string) (*model.Product, error) {
	product, err := s.repo.GetProductDetail(sku)
	if err != nil || product == nil {
		return nil, fmt.Errorf("query fail:%s", err)
	}
	return product, nil
}
