package controller

import (
	"bubble/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCateGoryItem(c *gin.Context) {
	var category models.Category
	c.BindJSON(&category)
	err := models.CreateCategoryItem(&category)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "添加成功",
			"data":   category,
		})
	}
}

func GetCategoryList(c *gin.Context) {
	categoryList, err := models.GetCategoryList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "获取Category列表成功",
			"data":   categoryList,
		})
	}
}

func GetOneCategory(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"err": "id未正确获取",
		})
	}
	category, err := models.GetCategoryItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取类别成功",
			"data": &category,
		})
	}
}

func UpdateCategoryItem(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"err": "id不存在",
		})
		return
	}
	category, err := models.GetCategoryItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.BindJSON(&category)
	if err := models.UpdateCategoryItem(category); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"msg":      "类别更新成功",
			"category": category,
		})
	}
}

func DeleteCategoryItem(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"error": "无效的id",
		})
	}
	err := models.DeleteCategoryItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "类别删除成功",
			"id":   id,
		})
	}
}
