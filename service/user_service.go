package service

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"bubble/global"
	"bubble/models"
	"bubble/pkg/convert"
	"bubble/pkg/errcode"

	"bubble/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

const TokenKey = "project manager user secret"

func GetToken(user *models.User) *models.JwtToken {

	info := make(map[string]interface{})
	now := time.Now()
	info["userId"] = user.ID
	info["exp"] = now.Add(time.Hour * 12).Unix() // 1 小时过期
	info["iat"] = now.Unix()
	tokenString := utils.CreateToken(TokenKey, info)
	return &models.JwtToken{
		Token: tokenString,
	}
}

func CheckToken(token string, user *models.User) (*models.UserInfo, *models.User, error) {
	info, ok := utils.ParseToken(token, TokenKey)
	infoMap := info.(jwt.MapClaims)
	if ok {
		expTime := infoMap["exp"].(float64)
		if float64(time.Now().Unix()) >= expTime {
			return nil, nil, fmt.Errorf("%s", "token已过期")
		} else {
			//根据id查询用户信息，并将用户信息系返回
			newUser, err := models.GetUserItem(strconv.Itoa(int(user.ID)))
			if err != nil {
				log.Println("checktoken获取用户信息失败", err)
			}
			info := &models.UserInfo{
				UserID:   int(newUser.ID),
				UserName: newUser.UserName,
				Phone:    newUser.Phone,
				Email:    newUser.Email,
			}
			return info, newUser, nil
		}
	} else {
		return nil, nil, fmt.Errorf("%s", "token无效")
	}
}

// ValidUsername 验证用户
func ValidUsername(username string) error {
	// 检测用户是否合规
	if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 12 {
		return errcode.UsernameLengthLimit
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		return errcode.UsernameCharLimit
	}

	// 重复检查
	user, _ := models.GetUserByUsername(username)

	if user != nil && user.ID > 0 {
		return errcode.UsernameHasExisted
	}

	return nil
}

// CheckPassword 密码检查
func CheckPassword(password string) error {
	// 检测用户是否合规
	if utf8.RuneCountInString(password) < 6 || utf8.RuneCountInString(password) > 16 {
		return errcode.PasswordLengthLimit
	}

	return nil
}

func ValidPassword(dbPassword, password, salt string) bool {
	return strings.Compare(dbPassword, utils.EncodeMD5(utils.EncodeMD5(password)+salt)) == 0
}

// EncryptPasswordAndSalt 密码加密&生成salt
func EncryptPasswordAndSalt(password string) (string, string) {
	salt := uuid.Must(uuid.NewV4()).String()[:8]
	password = utils.EncodeMD5(utils.EncodeMD5(password) + salt)

	return password, salt
}

// Register 用户注册
func Register(username, password string) (*models.User, error) {
	password, salt := EncryptPasswordAndSalt(password)

	user := &models.User{
		NickName: username,
		UserName: username,
		PassWord: password,
		Avatar:   GetRandomAvatar(),
		Salt:     salt,
		Status:   models.UserStatusNormal,
	}

	user, err := models.NewCreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const LOGIN_ERR_KEY = "UserLoginErr"
const MAX_LOGIN_ERR_TIMES = 10

// LoginAction 用户登录认证
func LoginAction(ctx *gin.Context, param *models.LoginRequest) (*models.User, error) {
	user, err := models.GetUserByUsername(param.Username)
	if err != nil {
		return nil, errcode.UnauthorizedAuthNotExist
	}

	if user != nil && user.ID > 0 {
		if errTimes, err := global.Redis.Get(ctx, fmt.Sprintf("%s:%d", LOGIN_ERR_KEY, user.ID)).Result(); err == nil {
			if convert.StrTo(errTimes).MustInt() >= MAX_LOGIN_ERR_TIMES {
				return nil, errcode.TooManyLoginError
			}
		}

		// 对比密码是否正确
		if ValidPassword(user.PassWord, param.Password, user.Salt) {

			if user.Status == models.UserStatusClosed {
				return nil, errcode.UserHasBeenBanned
			}

			// 清空登录计数
			global.Redis.Del(ctx, fmt.Sprintf("%s:%d", LOGIN_ERR_KEY, user.ID))
			return user, nil
		}

		// 登录错误计数
		_, err = global.Redis.Incr(ctx, fmt.Sprintf("%s:%d", LOGIN_ERR_KEY, user.ID)).Result()
		if err == nil {
			global.Redis.Expire(ctx, fmt.Sprintf("%s:%d", LOGIN_ERR_KEY, user.ID), time.Hour).Result()
		}

		return nil, errcode.UnauthorizedAuthFailed
	}

	return nil, errcode.UnauthorizedAuthNotExist
}

func UpdateUserInfo(user *models.User) error {
	return models.UpdateUserItem(user)
}

// GetUserInfo 获取用户信息
func GetUserInfo(param *models.AuthRequest) (*models.User, error) {
	user, err := models.GetUserByUsername(param.Username)

	if err != nil {
		return nil, err
	}

	if user != nil && user.ID > 0 {
		return user, nil
	}

	return nil, errcode.UnauthorizedAuthNotExist
}
