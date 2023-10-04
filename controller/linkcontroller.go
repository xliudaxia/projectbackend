package controller

import (
	"bubble/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateLinkItem(c *gin.Context) {
	var link models.Link
	c.BindJSON(&link)
	err := models.CreateLinkItem(&link)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "添加link成功",
			"data":   link,
		})
	}
}

func GetLinkList(c *gin.Context) {
	linkList, err := models.GetLinkList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "获取link列表成功",
			"data":   linkList,
		})
	}
}

func GetOneLink(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"err": "id未正确获取",
		})
	}
	link, err := models.GetLinkItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取link成功",
			"data": &link,
		})
	}
}

func UpdateLinkItem(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"err": "id不存在",
		})
		return
	}
	link, err := models.GetLinkItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.BindJSON(&link)
	if err := models.UpdateLinkItem(link); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "link更新成功",
			"link": link,
		})
	}
}

func DeleteLinkItem(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"error": "无效的id",
		})
	}
	err := models.DeleteLinkItem(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "link删除成功",
			"id":   id,
		})
	}
}

// 获取指定用户的link列表
func GetUserLinkList(c *gin.Context) {
	// 从中间件获取用户id
	userID, ok := c.Get("UID")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"error": "用户信息获取异常",
		})
	}

	linkList, err := models.GetUserLinkList(strconv.Itoa(userID.(int)))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"msg":     "获取用户link列表成功",
			"user_id": userID,
			"data":    linkList,
		})
	}
}
