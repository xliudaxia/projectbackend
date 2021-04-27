package service

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"bubble/models"

	"bubble/myutils"

	"github.com/dgrijalva/jwt-go"
)

const TokenKey = "project manager user secret"

func GetToken(user *models.User) *models.JwtToken {

	info := make(map[string]interface{})
	now := time.Now()
	info["userId"] = user.ID
	info["exp"] = now.Add(time.Hour * 12).Unix() // 1 小时过期
	info["iat"] = now.Unix()
	tokenString := myutils.CreateToken(TokenKey, info)
	return &models.JwtToken{
		Token: tokenString,
	}
}

func CheckToken(token string, user *models.User) (*models.UserInfo, error) {
	info, ok := myutils.ParseToken(token, TokenKey)
	infoMap := info.(jwt.MapClaims)
	if ok {
		expTime := infoMap["exp"].(float64)
		if float64(time.Now().Unix()) >= expTime {
			return nil, fmt.Errorf("%s", "token已过期")
		} else {
			//根据id查询用户信息，并将用户信息系返回
			newUser, err := models.GetUserItem(strconv.Itoa(user.ID))
			if err != nil {
				log.Println("checktoken获取用户信息失败", err)
			}
			info := &models.UserInfo{
				UserID:    newUser.ID,
				UserName:  newUser.UserName,
				Phone:     newUser.Phone,
				WorkCode:  newUser.WorkCode,
				CompanyID: newUser.CompanyID,
				Email:     newUser.Email,
			}
			return info, nil
		}
	} else {
		return nil, fmt.Errorf("%s", "token无效")
	}
}
