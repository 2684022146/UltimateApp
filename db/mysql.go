package db

import (
	"demo01/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
)

func InitMySQL(mysqlCfg config.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		mysqlCfg.Username,
		mysqlCfg.Password,
		mysqlCfg.Host,
		mysqlCfg.Port,
		mysqlCfg.DBname,
		mysqlCfg.DSNParams,
	)
	//  配置GORM日志（开发环境显示SQL，生产环境关闭）
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // 日志输出
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值
			LogLevel:      logger.Info, // 日志级别（dev:Info, prod:Error）
			Colorful:      true,        // 彩色输出
		},
	)

	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("db dsn fail:%s", err)
	}
	conn, _ := dbConn.DB()
	conn.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	conn.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	conn.SetConnMaxLifetime(func(s string) time.Duration {
		connMaxLifetime, _ := time.ParseDuration(s)
		return connMaxLifetime
	}(mysqlCfg.ConnMaxLifetime))
	if err = conn.Ping(); err != nil {
		panic(fmt.Sprintf("mysql connect fail:%s", err))
	}
	db = dbConn
	return db, nil
}
func GetDB() (*gorm.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("MySQL未初始化，请先调用InitMySQL()")
	}
	return db, nil
}
func CloseDB() error {
	if db != nil {

		dbConn, err := db.DB()
		if err != nil {
			return fmt.Errorf("获取底层sql.DB失败：%w", err)
		}
		log.Println("开始关闭MySQL连接池...")
		err = dbConn.Close()
		if err != nil {
			return fmt.Errorf("关闭MySQL连接池失败：%w", err)
		}
		log.Println("MySQL连接池已关闭")
		db = nil // 重置，防止重复关闭
	}
	return nil
}
