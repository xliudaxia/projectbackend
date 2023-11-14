package main

import (
	"bubble/controller"
	"bubble/middleware"
	"bubble/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	//设置前端打包目录的访问
	utils.UseSession(r)
	r.Static("/static", "static")
	// 加载静态文件目录
	// r.LoadHTMLGlob("templates/*")
	r.GET("/", controller.IndexHandler)

	v1Group := r.Group("v1")
	{
		// 用户登录
		v1Group.POST("/user/login", controller.NewUserLogin)
		// 用户注册
		v1Group.POST("/user/register", controller.NewRegister)
		v1Group.POST("/user/logout", controller.UserLogout)
		/* ******************验证码相关接口******************* */
		v1Group.GET("/captcha", controller.GetCaptcha)
		v1Group.POST("/captcha/verify", controller.VerifyCaptcha)
		/* ******************电话簿相关接口******************* */
		v1Group.POST("/phone", controller.CreatePhoneItem)
		v1Group.GET("/phone", controller.GetPhoneList)
		v1Group.GET("/phone/:id", controller.GetOnePhone)
		v1Group.PUT("/phone/:id", controller.UpdatePhoneItem)
		v1Group.DELETE("/phone/:id", controller.DeletePhoneItem)
		/* ******************书签目录相关接口******************* */
		v1Group.GET("/categoryList", controller.GetCategoryList)
		v1Group.POST("/category", controller.CreateCateGoryItem)
		v1Group.DELETE("/category/:id", controller.DeleteCategoryItem)
		v1Group.PUT("/category/:id", controller.UpdateCategoryItem)
		v1Group.GET("/category/:id", controller.GetOneCategory)
		// 获取广场流
		v1Group.GET("/posts", controller.GetPostList)
		// 获取指定用户的post列表
		v1Group.GET("/user/posts", controller.GetUserPosts)
		// 获取POST详情
		v1Group.GET("/post", controller.GetPost)

	}
	authApi := v1Group.Group("/").Use(middleware.JwtAuth())
	{
		// 获取用户信息
		authApi.GET("/user/info", controller.GetUserInfo)
		// 更改密码
		authApi.POST("/user/password", controller.ChangeUserPassword)
		// 更改昵称
		authApi.POST("/user/nickname", controller.ChangeNickName)
	}
	oldAuthApi := v1Group.Group("/").Use(middleware.UserAuthorize)
	{
		/* ******************书签相关接口******************* */
		oldAuthApi.GET("/linkList", controller.GetLinkList)
		oldAuthApi.POST("/link", controller.CreateLinkItem)
		oldAuthApi.DELETE("/link/:id", controller.DeleteLinkItem)
		oldAuthApi.PUT("/link/:id", controller.UpdateLinkItem)
		oldAuthApi.GET("/userLinkList", controller.GetUserLinkList)
		oldAuthApi.GET("/link/:id", controller.GetOneLink)
		/* ******************TODO相关接口********************/
		//添加todo记录
		oldAuthApi.POST("/todo", controller.CreateTodoItem)
		//查询全部todo数据
		oldAuthApi.GET("/todo", controller.GetTodoList)
		//根据ID查询单个todo记录
		oldAuthApi.GET("/todo/:id", controller.GetOneTodo)
		//根据ID修改TODO记录
		oldAuthApi.PUT("/todo/:id", controller.UpdateTodoItem)
		//根据ID删除某条记录
		oldAuthApi.DELETE("/todo/:id", controller.DeleteTodoItem)
		/* ******************用户相关接口******************* */
		oldAuthApi.POST("/user", controller.CreateUser)
		oldAuthApi.GET("/userlist", controller.GetUserList)
		oldAuthApi.GET("/currentUser", controller.CurrentUser)
		oldAuthApi.PUT("/user/:id", controller.UpdateUser)
		/* ******************项目相关接口******************* */
		oldAuthApi.GET("/projectlist", controller.GetProductList)
		oldAuthApi.POST("/project", controller.CreateProductItem)
		oldAuthApi.DELETE("/project", controller.DeleteProjectItem)
		oldAuthApi.PUT("/project", controller.UpdateProjectItem)
		oldAuthApi.GET("/queryproject", controller.QueryProjectItem)
		/* ******************Twitter相关接口******************* */
		// 新增POST方法
		oldAuthApi.POST("/post", controller.CreatePost)
		// 删除POST方法
		oldAuthApi.DELETE("/post", controller.DeletePost)
	}

	r.Run(":9090")
}
