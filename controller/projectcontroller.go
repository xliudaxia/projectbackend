package controller

import (
	"bubble/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var project models.Project

/*
CreateProductItem 创建一个项目
*/
func CreateProductItem(c *gin.Context) {
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
