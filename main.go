package main

import (
	"log"
	"net/http"
	"webdemo/config"
	"webdemo/db"
	"webdemo/route"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.InitConfig()
	dbConn, _ := db.InitMySQL(cfg.Mysql)
	log.Println("MySQL初始化成功")
	gin.SetMode(cfg.Server.Mode)
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
