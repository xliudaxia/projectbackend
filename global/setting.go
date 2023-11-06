package global

import (
	"time"
)

type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

var (
	JWTSetting *JWTSettingS
)
