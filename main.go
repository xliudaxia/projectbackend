package main

import (
	"bubble/controller"
	"bubble/dao"
	"bubble/models"

	"github.com/gin-gonic/gin"
)

func main() {
	err := dao.InitMysql()
	if err != nil {
		panic(err)
	}
	// defer dao.DB.close()
	dao.DB.AutoMigrate(&models.Todo{})
	dao.DB.AutoMigrate(&models.User{})
	dao.DB.AutoMigrate(&models.Project{})
	dao.DB.AutoMigrate(&models.Phone{})

	r := gin.Default()
	//设置前端打包目录的访问
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", controller.IndexHandler)

	v1Group := r.Group("v1")
	{
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
		//添加用户记录
		v1Group.POST("/user", controller.CreateUser)
		v1Group.GET("/userlist", controller.GetUserList)
		v1Group.GET("/currentUser", controller.CurrentUser)
		v1Group.PUT("/user/:id", controller.UpdateUser)
		v1Group.POST("/user/login", controller.UserLogin)
		v1Group.POST("/user/register", controller.UserRegister)
		/* ******************项目相关接口******************* */
		v1Group.GET("/projectlist", controller.GetProductList)
		v1Group.POST("/project", controller.CreateProductItem)
		v1Group.DELETE("/project", controller.DeleteProjectItem)
		v1Group.PUT("/project", controller.UpdateProjectItem)
		v1Group.GET("/queryproject", controller.QueryProjectItem)

		/* ******************电话簿相关接口******************* */
		//添加Phone记录
		v1Group.POST("/phone", controller.CreatePhoneItem)
		//查询全部Phone数据
		v1Group.GET("/phone", controller.GetPhoneList)
		//根据ID查询单个Phone记录
		v1Group.GET("/phone/:id", controller.GetOnePhone)
		//根据ID修改Phone记录
		v1Group.PUT("/phone/:id", controller.UpdatePhoneItem)
		//根据ID删除某条记录
		v1Group.DELETE("/phone/:id", controller.DeletePhoneItem)
	}
	r.Run(":9090")
}
