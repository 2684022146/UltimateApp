package route

import (
	"webdemo/controller"
	"webdemo/middleware"
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

	//登陆相关
	loginRepo := repository.NewLoginRepository(db)
	loginService := service.NewLoginService(loginRepo)
	loginController := controller.NewLoginController(loginService)
	r.POST("/login", loginController.Login)
	r.POST("/regist", loginController.Regist)
	//认证路由组
	authGroup := r.Group("/api")
	authGroup.Use(middleware.AuthMiddleware())
	//地址管理路由
	settingsRepo := repository.NewSettingsRepository(db)
	settingsService := service.NewSettingsService(settingsRepo)
	settingsController := controller.NewSettingsController(settingsService)
	consigneeGroup := authGroup.Group("/settings")
	consigneeGroup.Use(middleware.RequireRole())
	{
		//新建地址
		consigneeGroup.POST("/address", settingsController.CreateAddress)
		//地址列表
		consigneeGroup.GET("/address/list", settingsController.AddressList)

	}
	// // 认证路由组
	// authGroup := r.Group("/api")
	// authGroup.Use(middleware.AuthMiddleware())
	// {
	// 	// 退出登录（所有认证用户可访问）
	// 	authGroup.POST("/logout", settingController.Logout) // 退出登录

	// 	// 地址管理（仅收货人可访问）
	// 	consigneeGroup := authGroup.Group("/settings")
	// 	consigneeGroup.Use(middleware.RequireConsignee())
	// 	{
	// 		consigneeGroup.POST("/address", settingController.CreateAddress)                // 创建地址
	// 		consigneeGroup.PUT("/address", settingController.UpdateAddress)                 // 更新地址
	// 		consigneeGroup.DELETE("/address", settingController.DeleteAddress)              // 删除地址
	// 		consigneeGroup.GET("/address/:id", settingController.GetAddressByID)            // 获取地址详情
	// 		consigneeGroup.GET("/addresses", settingController.GetAddressList)              // 获取地址列表
	// 		consigneeGroup.PUT("/address/:id/default", settingController.SetDefaultAddress) // 设置默认地址
	// 	}
	// }

	return r
}
