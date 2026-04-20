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
		//更新地址
		consigneeGroup.PUT("address", settingsController.UpdateAddress)
		//地址详情
		consigneeGroup.GET("/address", settingsController.AddressDetail)
		//删除地址
		consigneeGroup.DELETE("/address", settingsController.DeleteAddress)
		//设为默认地址
		consigneeGroup.PUT("/address/default", settingsController.SetDefault)

	}
	//订单管理路由
	ordersRepo := repository.NewOrdersRepository(db)
	ordersService := service.NewOrdersService(ordersRepo)
	ordersController := controller.NewOrdersController(ordersService)

	// 客户订单路由组
	consigneeOrdersGroup := authGroup.Group("/orders")
	consigneeOrdersGroup.Use(middleware.RequireRole())
	{

		//创建订单
		consigneeOrdersGroup.POST("/create", ordersController.CreateOrder)
	}
	senderOrderGroup := consigneeOrdersGroup.Group("/sender")
	senderOrderGroup.Use(middleware.RequireRole())
	{
		//获取寄件人已完成订单 我寄出的已完成订单
		senderOrderGroup.GET("/finished", ordersController.SenderFinishedOrder)
		//获取我在寄出的途中订单 我寄出的在途中订单
		senderOrderGroup.GET("/delivering", ordersController.SenderInTransitOrder)
		//获取我在寄出的待接单订单 我寄出的待接单订单
		senderOrderGroup.GET("/waiting", ordersController.SenderWaitingOrder)
		//取消订单
		senderOrderGroup.POST("/cancel", ordersController.SenderCencelOrder)

	}
	receiverOrderGroup := consigneeOrdersGroup.Group("/receiver")
	receiverOrderGroup.Use(middleware.RequireRole())
	{
		//获取收件人已完成订单
		receiverOrderGroup.GET("/finished", ordersController.ReceiverFinishedOrder)
		//获取收件人途中订单
		receiverOrderGroup.GET("/delivering", ordersController.ReceiverInTransitOrder)
	}
	riderGroup := authGroup.Group("/rider")
	//暂时注释掉权限中间件，以便测试骑手接单的业务逻辑
	//riderGroup.Use(middleware.RequireRole())
	{
		//骑手待接单订单
		riderGroup.GET("/waiting", ordersController.RiderOrderList)
		//骑手接单
		riderGroup.POST("/accept", ordersController.AcceptOrder)
		//开始配送
		riderGroup.POST("/start", ordersController.StartDelivery)
		//上传位置
		riderGroup.POST("/location", ordersController.UploadLocation)
		//完成订单
		riderGroup.POST("/complete", ordersController.CompleteOrder)
	}

	return r
}
