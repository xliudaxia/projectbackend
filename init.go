package main

import (
	"bubble/dao"
	"bubble/global"
	"bubble/models"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

// 初始化jwt和redis
func init() {

	// 加载环境变量
	err := godotenv.Load(".env.local")
	if err != nil {
		panic(err)
	}

	// 初始化数据库
	err = dao.InitMysql()
	if err != nil {
		panic(err)
	}

	// defer dao.DB.close()
	dao.DB.AutoMigrate(&models.Todo{}, &models.User{}, &models.Project{}, &models.Phone{}, &models.Category{}, &models.Link{}, &models.Post{}, &models.PostContent{})

	global.JWTSetting = &global.JWTSettingS{
		Secret: os.Getenv("JWT_SECRET"),
		Issuer: os.Getenv("JWT_ISSUER"),
		Expire: 86400,
	}

	fmt.Println("redis debug", os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWD"))

	global.Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWD"),
	})
	global.JWTSetting.Expire *= time.Second
}
