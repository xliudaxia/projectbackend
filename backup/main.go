package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

//创建全局DB对象
var (
	DB *gorm.DB
)

//连接数据库方法
func initMysql() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}

func main() {
	err := initMysql()
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	//设置前端打包目录的访问
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1Group := r.Group("v1")
	{
		//添加todo记录
		v1Group.POST("/todo", func(c *gin.Context) {
			var todo Todo
			c.BindJSON(&todo)
			if err = DB.Create(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"msg":    "添加成功",
					"data":   todo,
				})
			}
		})
		//查询全部todo数据
		v1Group.GET("/todo", func(c *gin.Context) {
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})
		//根据ID查询单个todo记录
		v1Group.GET("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"err": "id未正确获取",
				})
			}
			var todo Todo
			if err = DB.Where("id = ?", id).First(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"err": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "数据查询成功",
					"data": &todo,
				})
			}
		})
		//根据ID修改TODO记录
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"err": "id不存在",
				})
			}
			var todo Todo
			if err = DB.Where("id = ?", id).First(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			}
			c.BindJSON(&todo)
			// fmt.Println("操作后todo,%#v", todo)
			if err = DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {

				c.JSON(http.StatusOK, todo)
			}
		})
		//根据ID删除某条记录
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"error": "无效的id",
				})
			}
			if err = DB.Where("id=?", id).Delete(Todo{}).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					id: "deleted",
				})
			}
		})
	}
	r.Run(":9090")
}
