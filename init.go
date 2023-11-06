package main

import (
	"bubble/global"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// 初始化jwt和redis
func init() {

	global.JWTSetting = &global.JWTSettingS{
		Secret: os.Getenv("JWT_SECRET"),
		Issuer: os.Getenv("JWT_ISSUER"),
		Expire: 86400,
	}

	global.Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWD"),
	})
	global.JWTSetting.Expire *= time.Second
}
