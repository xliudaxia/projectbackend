package controller

import (
	"bubble/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jonluo94/cool/gintool"
)

/*
CreateProductItem 创建一个项目
*/
func CreateProductItem(c *gin.Context) {
	var project models.Project
	c.BindJSON(&project)
	err := models.CreateProjectItem(&project)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "添加成功",
			"data":   project,
		})
	}
}

/*
GetProductList 获取项目列表
*/
func GetProductList(c *gin.Context) {
	projectList, err := models.GetProjectList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, projectList)
	}
}

type DelProjectArgs struct {
	ID int `json:"id"`
}

/*
DeleteProjectItem 删除项目记录接口
*/
func DeleteProjectItem(c *gin.Context) {
	var tempproj DelProjectArgs
	c.BindJSON(&tempproj)
	fmt.Println("获取到的todo结果%#v", tempproj)
	err := models.DeleteProjectItem(strconv.Itoa(tempproj.ID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "删除成功",
		})
	}
}

/*
UpdateProjectItem 更新单条项目记录接口
*/
func UpdateProjectItem(c *gin.Context) {
	var tempproj models.Project
	c.ShouldBindBodyWith(&tempproj, binding.JSON)
	// 通过project_id在库中查询记录
	currentproject, err := models.GetProjectItem(strconv.Itoa(tempproj.ID))
	if err != nil {
		gintool.ResultFail(c, err)
		return
	}
	//将返回的数据绑定给currentproject
	c.ShouldBindBodyWith(&currentproject, binding.JSON)
	//执行更新操作
	if err := models.UpdateProjectItem(currentproject); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "项目信息更新成功！",
		})
	}
}

/*
QueryProjectItem 通过项目名称或项目介绍查询项目
*/
func QueryProjectItem(c *gin.Context) {
	queryStr := c.Query("keyword")
	if projectList, err := models.QueryProjectByRule(queryStr); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":      http.StatusOK,
			"msg":         "项目查询结果返回成功！",
			"projectlist": projectList,
		})
	}

}
