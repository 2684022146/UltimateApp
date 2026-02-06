package route

import (
	"webdemo/controller"
	"webdemo/repository"
	"webdemo/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	//跨域
	r.Use(cors.New(cors.Config{
		// 允许的前端来源（生产环境替换为你的前端域名，如 http://localhost:5500）
		AllowOrigins: []string{"*"},
		// 允许的请求方法（GET/POST/PUT/DELETE 等）
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 允许的请求头（Content-Type、Authorization 等）
		AllowHeaders: []string{"*"},
		// 允许前端携带 Cookie/Token（如需认证则开启）
		AllowCredentials: true,
		// 预检请求缓存时间（秒），减少 OPTIONS 请求次数
		MaxAge: 86400,
	}))

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productController := controller.NewProductController(productService)
	productGroup := r.Group("/product")
	{
		productGroup.GET("/list", productController.GetProductList)
		productGroup.GET("/detail", productController.ProductDetail)
	}
	loginRepo := repository.NewLoginRepository(db)
	loginService := service.NewLoginService(loginRepo)
	loginController := controller.NewLoginController(loginService)
	r.POST("/login", loginController.Login)
	r.POST("/regist", loginController.Regist)
	return r
}
