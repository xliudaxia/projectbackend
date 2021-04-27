package controller

import (
	"bubble/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePhoneItem(c *gin.Context) {
	var phone models.Phone
	c.ShouldBindJSON(&phone)
	err := models.CreatePhoneItem(&phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "添加成功",
			"data":   phone,
		})
	}
}

func GetPhoneList(c *gin.Context) {
	phoneList, err := models.GetPhoneList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "获取电话簿列表成功",
			"data":   phoneList,
		})
	}
}

func GetOnePhone(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"err": "id未正确获取",
		})
	}
	phone, err := models.GetPhoneItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "数据查询成功",
			"data": &phone,
		})
	}
}

func UpdatePhoneItem(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"err": "id不存在",
		})
		return
	}
	phone, err := models.GetPhoneItem(id)
	fmt.Println("获取到的phone结果%#v", phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.ShouldBindJSON(&phone)
	if err := models.UpdatePhoneItem(phone); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  500,
			"msg":   "更新联系人失败",
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新联系人信息成功",
			"data": &phone,
		})
	}
}

func DeletePhoneItem(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"error": "无效的id",
		})
	}
	err := models.DeletePhoneItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  500,
			"msg":   "联系人删除失败，请稍后再试",
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "联系人删除成功",
			"id":   id,
		})
	}
}
