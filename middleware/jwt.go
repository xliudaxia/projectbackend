package middleware

import (
	"bubble/models"
	"bubble/myutils"
	"bubble/service"
	"net/http"

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
	user, err := service.CheckToken(token, &models.User{ID: session.(int)})

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
		c.Next()
	}
}
