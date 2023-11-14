package controller

import (
	"bubble/models"
	"bubble/pkg/app"
	"bubble/pkg/errcode"
	"bubble/service"
	"bubble/utils"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

var user models.User

// CreateUser 创建新的用户
func CreateUser(c *gin.Context) {
	c.BindJSON(&user)
	err := models.CreateUserItem(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "添加成功",
			"data":   user,
		})
	}

}

// GetUserList ：获取用户列表
func GetUserList(c *gin.Context) {
	userList, err := models.GetUserList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, userList)
	}
}

// UpdateUser 更新用户记录
func UpdateUser(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"err": "id不存在",
		})
		return
	}
	user, err := models.GetUserItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	}
	c.BindJSON(&user)
	if err = models.UpdateUserItem(user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "更新成功",
			"user":    user,
		})
	}
}

// UserLogin 用户登录接口
func UserLogin(c *gin.Context) {
	/*
	   Token校验实现步骤：
	   1、登录系统根据相关信息生成token  √
	   2、生成完token后，登录时将token返回并在cookies中保存   √×
	   3、用户根据token信息获取当前用户信息   √
	   4、登录后请求的接口需要增加token校验，不符合则经过中间件进行拦截
	*/
	userdata, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("获取用户信息失败! err-->", err)
		c.String(http.StatusOK, err.Error())
		return
	}
	userreq := models.User{}
	err = json.Unmarshal(userdata, &userreq)
	if err != nil {
		log.Println("用户绑定到结构体失败! err-->", err)
		c.String(http.StatusOK, err.Error())
		return
	}
	//判断用户名密码是否为空
	if userreq.UserName == "" || userreq.PassWord == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "表单数据不完整，请重试！",
		})
		return
	}
	//在数据库中查找当前登录用户的用户名
	queryuser, err := models.QueryLoginUser(userreq.UserName, "user_name")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":           "error",
			"type":             "account",
			"currentAuthority": "guest",
			"message":          "用户名或密码错误，请重试！",
		})
		return
	}
	//验证当前登录用户的密码
	if queryuser.PassWord != userreq.PassWord {
		c.JSON(http.StatusOK, gin.H{
			"status":           "error",
			"type":             "account",
			"currentAuthority": "guest",
			"message":          "用户名或密码错误，请重试！",
		})
		return
	}
	//执行token生成操作
	token := service.GetToken(queryuser)
	utils.SetSession(c, token.Token, queryuser.ID)

	c.JSON(http.StatusOK, gin.H{
		"status":           "ok",
		"type":             "account",
		"avatar":           "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
		"notifyCount":      12,
		"unreadCount":      11,
		"currentAuthority": "admin",
		"user":             queryuser,
		"token":            token.Token,
		"message":          "登录成功",
	})
}

func NewUserLogin(c *gin.Context) {
	param := models.LoginRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)

	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	user, err := service.LoginAction(c, &param)

	if err != nil {
		response.ToErrorResponse(err.(*errcode.Error))
		return
	}

	token, err := app.GenerateToken(user)

	// 兼容性操作
	utils.SetSession(c, token, user.ID)

	if err != nil {
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		return
	}

	response.ToResponse(gin.H{
		"token": token,
	})

}

// 用户退出接口
func UserLogout(c *gin.Context) {
	token := c.GetHeader("M-Token")
	utils.RemoveSession(c, token)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
		"msg":  "退出成功",
	})
}

func ChangeUserPassword(c *gin.Context) {
	param := models.ChangePasswordReq{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	user := &models.User{}
	if u, exists := c.Get("USER"); exists {
		user = u.(*models.User)
	}

	// 密码检查
	err := service.CheckPassword(param.Password)
	if err != nil {
		response.ToErrorResponse(err.(*errcode.Error))
		return
	}

	// 旧密码校验
	if !service.ValidPassword(user.PassWord, param.OldPassword, user.Salt) {
		response.ToErrorResponse(errcode.ErrorOldPassword)
		return
	}

	user.PassWord, user.Salt = service.EncryptPasswordAndSalt(param.Password)

	service.UpdateUserInfo(user)

	response.ToResponse(nil)
}

func ChangeNickName(c *gin.Context) {
	param := models.ChangeNicknameReq{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	user := &models.User{}
	if u, exists := c.Get("USER"); exists {
		user = u.(*models.User)
	}

	if utf8.RuneCountInString(param.Nickname) < 2 || utf8.RuneCountInString(param.Nickname) > 12 {
		response.ToErrorResponse(errcode.NicknameLengthLimit)
		return
	}

	user.NickName = param.Nickname
	service.UpdateUserInfo(user)
	response.ToResponse(nil)
}

// currentUser 获取用户信息接口
func CurrentUser(c *gin.Context) {
	token := c.GetHeader("M-Token")
	session := utils.GetSession(c, token)
	if nil == session {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  "token不存在"})
		return
	}
	user, _, err := service.CheckToken(token, &models.User{ID: session.(int64)})
	if err != nil {
		if err.Error() == "token已过期" || err.Error() == "token无效" {
			m := make(map[string]interface{})
			m["code"] = 2
			m["msg"] = err.Error()
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"data": nil,
				"msg":  "token已失效"})
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err.Error()})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	c.BindJSON(&user)
	//注册参数完整性校验
	if user.UserName == "" || user.Email == "" || user.Phone == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "参数不完整，请检查！",
		})
		return
	}
	//校验当前用户名是否已经注册
	tmpuser, err := models.RegisterUserCheck(user.UserName, "user_name")
	if err != nil {
		log.Println("查询当前用户是否已注册失败！")
		c.JSON(http.StatusOK, err.Error())
		return
	}
	if tmpuser.UserName != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "该用户名已被注册，请更换后重试！",
		})
		return
	}
	//校验当前注册用户邮箱是否重复
	tmpuser, err = models.RegisterUserCheck(user.Email, "email")
	if err != nil {
		log.Println("查询注册用户邮箱是否重复失败！")
		c.JSON(http.StatusOK, err.Error())
		return
	}
	if tmpuser.UserName != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "当前注册用户邮箱已被占用，请更换后重试！",
		})
		return
	}
	//校验当前注册用户手机是否重复
	tmpuser, err = models.RegisterUserCheck(strconv.Itoa(user.Phone), "phone")
	if err != nil {
		log.Println("查询注册用户手机是否重复失败！")
		c.JSON(http.StatusOK, err.Error())
		return
	}
	if tmpuser.UserName != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "当前注册用户手机已被占用，请更换后重试！",
		})
		return
	}
	//执行用户注册逻辑
	err = models.CreateUserItem(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  1,
			"error": err.Error(),
		})
	} else {
		user.PassWord = ""
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "用户注册成功",
			"data":   user,
		})
	}

}

func NewRegister(c *gin.Context) {
	param := models.RegisterRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	// 用户名检查
	err := service.ValidUsername(param.Username)
	if err != nil {
		response.ToErrorResponse(err.(*errcode.Error))
		return
	}

	// 密码检查
	err = service.CheckPassword(param.Password)
	if err != nil {
		response.ToErrorResponse(err.(*errcode.Error))
		return
	}

	user, err := service.Register(
		param.Username,
		param.Password,
	)

	if err != nil {
		response.ToErrorResponse(errcode.UserRegisterFailed)
		return
	}

	response.ToResponse(gin.H{
		"id":       user.ID,
		"username": user.UserName,
	})

}

func GetUserInfo(c *gin.Context) {
	param := models.AuthRequest{}
	response := app.NewResponse(c)

	if username, exists := c.Get("USERNAME"); exists {
		param.Username = username.(string)
	}

	user, err := service.GetUserInfo(&param)

	if err != nil {
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}

	phone := ""

	if strconv.Itoa(user.Phone) != "" && len(strconv.Itoa(user.Phone)) == 11 {
		phone = strconv.Itoa(user.Phone)[0:3] + "****" + strconv.Itoa(user.Phone)[7:]
	}

	response.ToResponse(gin.H{
		"id":       user.ID,
		"nickname": user.NickName,
		"username": user.UserName,
		"status":   user.Status,
		"avatar":   user.Avatar,
		"phone":    phone,
		"is_admin": user.IsAdmin,
	})
}
