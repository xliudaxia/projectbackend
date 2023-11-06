package middleware

import (
	"bubble/models"
	"bubble/myutils"
	"bubble/pkg/app"
	"bubble/pkg/errcode"
	"bubble/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// UserAuthorize 用户权限校验中间件
func UserAuthorize(c *gin.Context) {
	var token string
	var err error
	m := make(map[string]interface{})
	m["code"] = 2
	token = c.GetHeader("M-Token")
	if token == "" {
		token, err = c.Cookie("M-Token")
		if err != nil {
			m["msg"] = err.Error()
			c.JSON(http.StatusOK, m)
			c.Abort()
			return
		}
	}
	session := myutils.GetSession(c, token)
	if nil == session {
		m["msg"] = "token不存在"
		c.JSON(http.StatusOK, m)
		c.Abort()
		return
	}
	user, originUser, err := service.CheckToken(token, &models.User{ID: session.(int64)})

	if err != nil {
		if err.Error() == "token已过期" || err.Error() == "token无效" {
			m["msg"] = err.Error()
			c.JSON(http.StatusOK, m)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"data": nil,
				"msg":  err.Error(),
			})
		}
		c.Abort()
		return
	} else {
		c.Set("UID", user.UserID)
		c.Set("USERNAME", user.UserName)
		c.Set("USER", originUser)
		c.Next()
	}
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			ecode = errcode.Success
		)
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("Authorization")

			// 验证前端传过来的token格式，不为空，开头为Bearer
			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				response := app.NewResponse(c)
				response.ToErrorResponse(errcode.UnauthorizedTokenError)
				c.Abort()
				return
			}

			// 验证通过，提取有效部分（除去Bearer)
			token = token[7:]
		}
		if token == "" {
			ecode = errcode.InvalidParams
		} else {
			claims, err := app.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ecode = errcode.UnauthorizedTokenTimeout
				default:
					ecode = errcode.UnauthorizedTokenError
				}
			} else {
				c.Set("UID", claims.UID)
				c.Set("USERNAME", claims.Username)

				// 加载用户信息
				user := &models.User{
					ID: claims.UID,
				}
				user, _ = models.GetUserItem(strconv.Itoa(int(claims.UID)))
				c.Set("USER", user)

				// 强制下线机制
				if ("paopao-api" + ":" + user.Salt) != claims.Issuer {
					ecode = errcode.UnauthorizedTokenTimeout
				}
			}
		}

		if ecode != errcode.Success {
			response := app.NewResponse(c)
			response.ToErrorResponse(ecode)
			c.Abort()
			return
		}

		c.Next()
	}
}
