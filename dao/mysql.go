package dao

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 创建全局DB对象
var (
	DB *gorm.DB
)

// InitMysql 连接数据库方法
func InitMysql() (err error) {
	// 数据库连接地址
	DBIpAdress := os.Getenv("DB_IP_ADRESSS")
	// 数据库用户名
	DBUserName := os.Getenv("DB_USERNAME")
	// 数据库密码
	DBPassword := os.Getenv("DB_PASSWORD")
	// 数据库端口号
	DBPort := os.Getenv("DB_PORT")
	dsn := DBUserName + ":" + DBPassword + "@tcp(" + DBIpAdress + ":" + DBPort + ")/projectmanager?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}
