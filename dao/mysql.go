package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//创建全局DB对象
var (
	DB *gorm.DB
)

//InitMysql 连接数据库方法
//TODO 上传前需脱敏
func InitMysql() (err error) {
	dsn := "root:m52191061@tcp(120.53.241.206:3306)/projectmanager?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}
