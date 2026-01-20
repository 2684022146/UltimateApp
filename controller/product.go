package controller

import (
	"demo01/response"
	"demo01/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService service.ProductService
}

func NewProductController(c service.ProductService) *ProductController {
	return &ProductController{
		productService: c,
	}
}

func (c *ProductController) GetProductList(ctx *gin.Context) {
	productList, err := c.productService.ProductList(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &response.BaseResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}
	var resp []*response.ProductResponse
	for _, v := range productList {
		productLD := &response.ProductJSONLD{
			Context: "https://Schema.org/",
			Type:    "Product",
			Name:    v.Name,
			Price:   v.Price,
		}
		resp = append(resp, &response.ProductResponse{
			Name:          v.Name,
			Price:         v.Price,
			ProductJSONLD: *productLD,
			Sku:           v.Sku,
		})
	}

	ctx.JSON(http.StatusOK, &response.BaseResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    resp,
	})
}
func (c *ProductController) ProductDetail(ctx *gin.Context) {
	sku := ctx.Query("sku")
	if strings.TrimSpace(sku) == "" {
		ctx.JSON(400, response.BaseResponse{
			Code:    400,
			Message: "参数错误",
		})
	}
	product, err := c.productService.ProductDetail(ctx, sku)
	if err != nil || product == nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Code:    400,
			Message: err.Error(),
		})
	}
	productLD := response.ProductJSONLD{
		Context: "https://Schema.org/",
		Type:    "Product",
		Name:    product.Name,
		Price:   product.Price,
	}
	resp := response.ProductResponse{
		Name:          product.Name,
		Price:         product.Price,
		Description:   product.Description,
		Sku:           product.Sku,
		ProductJSONLD: productLD,
	}
	ctx.JSON(http.StatusOK, response.BaseResponse{
		Code:    200,
		Message: "success",
		Data:    resp,
	})
}
