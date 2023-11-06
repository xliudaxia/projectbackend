package main

import (
	"bubble/controller"
	"bubble/dao"
	"bubble/middleware"
	"bubble/models"
	"bubble/myutils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

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

	r := gin.Default()
	//设置前端打包目录的访问
	myutils.UseSession(r)
	r.Static("/static", "static")
	// 加载静态文件目录
	// r.LoadHTMLGlob("templates/*")
	r.GET("/", controller.IndexHandler)

	v1Group := r.Group("v1")

	authApi := v1Group.Group("/").Use(middleware.JwtAuth())
	{
		// 获取用户信息
		authApi.GET("/user/info", controller.GetUserInfo)
		// 更改密码
		authApi.POST("/user/password", controller.ChangeUserPassword)
		// 更改昵称
		authApi.POST("/user/nickname", controller.ChangeNickName)
	}
	{
		// 用户登录
		v1Group.POST("/user/login", controller.NewUserLogin)
		// 用户注册
		v1Group.POST("/user/register", controller.NewRegister)
		v1Group.POST("/user/logout", controller.UserLogout)
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
		//权限校验
		v1Group.Use(middleware.UserAuthorize)
		/* ******************书签相关接口******************* */
		v1Group.GET("/linkList", controller.GetLinkList)
		v1Group.POST("/link", controller.CreateLinkItem)
		v1Group.DELETE("/link/:id", controller.DeleteLinkItem)
		v1Group.PUT("/link/:id", controller.UpdateLinkItem)
		v1Group.GET("/userLinkList", controller.GetUserLinkList)
		v1Group.GET("/link/:id", controller.GetOneLink)
		/* ******************TODO相关接口********************/
		//添加todo记录
		v1Group.POST("/todo", controller.CreateTodoItem)
		//查询全部todo数据
		v1Group.GET("/todo", controller.GetTodoList)
		//根据ID查询单个todo记录
		v1Group.GET("/todo/:id", controller.GetOneTodo)
		//根据ID修改TODO记录
		v1Group.PUT("/todo/:id", controller.UpdateTodoItem)
		//根据ID删除某条记录
		v1Group.DELETE("/todo/:id", controller.DeleteTodoItem)
		/* ******************用户相关接口******************* */
		v1Group.POST("/user", controller.CreateUser)
		v1Group.GET("/userlist", controller.GetUserList)
		v1Group.GET("/currentUser", controller.CurrentUser)
		v1Group.PUT("/user/:id", controller.UpdateUser)
		/* ******************项目相关接口******************* */
		v1Group.GET("/projectlist", controller.GetProductList)
		v1Group.POST("/project", controller.CreateProductItem)
		v1Group.DELETE("/project", controller.DeleteProjectItem)
		v1Group.PUT("/project", controller.UpdateProjectItem)
		v1Group.GET("/queryproject", controller.QueryProjectItem)
		/* ******************Twitter相关接口******************* */
		// 新增POST方法
		v1Group.POST("/post", controller.CreatePost)
		// 删除POST方法
		v1Group.DELETE("/post", controller.DeletePost)
	}

	r.Run(":9090")
}
