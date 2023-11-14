package controller

import (
	"bubble/models"
	"bubble/pkg/app"
	"bubble/pkg/convert"
	"bubble/pkg/errcode"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	param := models.PostCreationReq{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}
	userID, _ := c.Get("UID")

	post := &models.Post{
		UserID: userID.(int64),
		Tags:   strings.Join(param.Tags, ","),
		IP:     "127.0.0.2",
		IPLoc:  "demo location",
	}
	postID, err := models.CreatePostItem(post)
	if err != nil {
		return
	}

	for _, item := range param.Contents {
		postContent := &models.PostContent{
			PostID:  int64(postID),
			UserID:  userID.(int64),
			Content: item.Content,
			Type:    item.Type,
			Sort:    item.Sort,
		}
		models.CreatePostContent(postContent)
	}
	response.ToResponse(post)
}

func DeletePost(c *gin.Context) {
	param := models.PostDelReq{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	user, _ := c.Get("USER")

	postFormated, err := models.GetPost(param.ID)

	if err != nil {
		response.ToErrorResponse(errcode.GetPostFailed)
		return
	}

	if postFormated.UserID != user.(*models.User).ID && !user.(*models.User).IsAdmin {
		response.ToErrorResponse(errcode.NoPermission)
		return
	}

	err = models.DeletePostItem(strconv.FormatInt(int64(param.ID), 10))

	if err != nil {
		return
	}
	response.ToResponse(nil)
}

func GetUserPosts(c *gin.Context) {
	response := app.NewResponse(c)
	username := c.Query("username")

	user, err := models.GetUserByUsername(username)
	if err != nil {
		response.ToErrorResponse(errcode.NoExistUsername)
		return
	}

	conditions := &models.ConditionsT{
		"user_id": user.ID,
		"ORDER":   "latest_replied_on DESC",
	}

	posts, err := models.GetPostList(&models.PostListReq{
		Conditions: conditions,
		Offset:     (app.GetPage(c) - 1) * app.GetPageSize(c),
		Limit:      app.GetPageSize(c),
	})

	if err != nil {
		response.ToErrorResponse(errcode.GetPostsFailed)
		return
	}

	totalRows, _ := models.GetPostCount(conditions)

	response.ToResponseList(posts, totalRows)
}

func GetPost(c *gin.Context) {
	postID := convert.StrTo(c.Query("id")).MustInt64()
	response := app.NewResponse(c)

	postFormated, err := models.GetPost(postID)

	if err != nil {
		response.ToErrorResponse(errcode.GetPostFailed)
		return
	}

	response.ToResponse(postFormated)
}

func GetPostList(c *gin.Context) {
	response := app.NewResponse(c)
	// 直接读库
	posts, err := models.GetPostList(&models.PostListReq{
		Conditions: &models.ConditionsT{
			"ORDER": "latest_replied_on DESC",
		},
		Offset: (app.GetPage(c) - 1) * app.GetPageSize(c),
		Limit:  app.GetPageSize(c),
	})
	if err != nil {
		response.ToErrorResponse(errcode.GetPostsFailed)
		return
	}
	totalRows, _ := models.GetPostCount(&models.ConditionsT{
		"ORDER": "latest_replied_on DESC",
	})

	response.ToResponseList(posts, totalRows)
}
