package controller

import (
	"bubble/global"
	"bubble/pkg/app"
	"bubble/pkg/errcode"
	"bubble/utils"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"
	"time"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type VerifyCaptchaReq struct {
	ImgCaptcha   string `json:"img_captcha" form:"img_captcha" binding:"required"`
	ImgCaptchaID string `json:"img_captcha_id" form:"img_captcha_id" binding:"required"`
}

func GetCaptcha(c *gin.Context) {
	cap := captcha.New()

	if err := cap.SetFont("static/comic.ttf"); err != nil {
		panic(err.Error())
	}

	cap.SetSize(160, 64)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{0, 0, 0, 255})
	cap.SetBkgColor(color.RGBA{218, 240, 228, 255})
	img, password := cap.Create(6, captcha.NUM)
	emptyBuff := bytes.NewBuffer(nil)
	_ = png.Encode(emptyBuff, img)

	key := utils.EncodeMD5(uuid.Must(uuid.NewV4()).String())

	// 五分钟有效期
	global.Redis.SetEX(c, "Captcha:"+key, password, time.Minute*5)

	response := app.NewResponse(c)
	response.ToResponse(gin.H{
		"id":   key,
		"b64s": "data:image/png;base64," + base64.StdEncoding.EncodeToString(emptyBuff.Bytes()),
	})
}

func VerifyCaptcha(c *gin.Context) {
	param := VerifyCaptchaReq{}

	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	fmt.Println(global.Redis.Get(c.Request.Context(), "Captcha:"+param.ImgCaptchaID))

	// 验证图片验证码
	if res, err := global.Redis.Get(c.Request.Context(), "Captcha:"+param.ImgCaptchaID).Result(); err != nil || res != param.ImgCaptcha {
		response.ToErrorResponse(errcode.ErrorCaptchaPassword)
		return
	}
	global.Redis.Del(c.Request.Context(), "Captcha:"+param.ImgCaptchaID).Result()

	response.ToResponse(nil)
}
