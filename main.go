package main

import (
	"log"
	"net/http"
	"webdemo/config"
	"webdemo/db"
	"webdemo/middleware"
	"webdemo/repository"
	"webdemo/route"
	"webdemo/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.InitConfig()
	dbConn, _ := db.InitMySQL(cfg.Mysql)
	log.Println("MySQL初始化成功")
	gin.SetMode(cfg.Server.Mode)
	permRepo := repository.NewPermissionRepository(dbConn)
	permService := service.NewPermissionService(permRepo)
	middleware.SetPermissionService(permService)
	r := route.InitRouter(dbConn)
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}
	go func() {
		log.Println("server start:", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	select {}
}
