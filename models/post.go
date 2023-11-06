package models

import (
	"bubble/dao"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostContentT int

type Post struct {
	gorm.Model
	UserID          int64  `json:"user_id"`
	CommentCount    int64  `json:"comment_count"`
	LatestRepliedOn int64  `json:"latest_replied_on"`
	Tags            string `json:"tags"`
	IP              string `json:"ip"`
	IPLoc           string `json:"ip_loc"`
}

type PostFormated struct {
	ID              int64                  `json:"id"`
	UserID          int64                  `json:"user_id"`
	User            *UserFormated          `json:"user"`
	Contents        []*PostContentFormated `json:"contents"`
	CommentCount    int64                  `json:"comment_count"`
	LatestRepliedOn int64                  `json:"latest_replied_on"`
	CreatedOn       time.Time              `json:"created_on"`
	ModifiedOn      time.Time              `json:"modified_on"`
	Tags            map[string]int8        `json:"tags"`
	IPLoc           string                 `json:"ip_loc"`
}

type PostContent struct {
	gorm.Model
	PostID  int64  `json:"post_id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
	Type    int    `json:"type"`
	Sort    int64  `json:"sort"`
}

type PostContentFormated struct {
	ID      int64        `json:"id"`
	PostID  int64        `json:"post_id"`
	Content string       `json:"content"`
	Type    PostContentT `json:"type"`
	Sort    int64        `json:"sort"`
}

type PostContentItem struct {
	Content string `json:"content"  binding:"required"`
	Type    int    `json:"type"  binding:"required"`
	Sort    int64  `json:"sort"  binding:"required"`
}

type PostCreationReq struct {
	Contents []*PostContentItem `json:"contents" binding:"required"`
	Tags     []string           `json:"tags" binding:"required"`
}

type PostDelReq struct {
	ID int64 `json:"id" binding:"required"`
}

type PostListReq struct {
	Conditions *ConditionsT
	Offset     int
	Limit      int
}

func (p *PostContent) Format() *PostContentFormated {
	if p == nil {
		return nil
	}
	return &PostContentFormated{
		ID:      int64(p.ID),
		PostID:  int64(p.ID),
		Content: p.Content,
		Type:    PostContentT(p.Type),
		Sort:    p.Sort,
	}
}

func (p *Post) Format() *PostFormated {
	if p != nil {
		tagsMap := map[string]int8{}
		for _, tag := range strings.Split(p.Tags, ",") {
			tagsMap[tag] = 1
		}
		return &PostFormated{
			ID:              int64(p.ID),
			UserID:          p.UserID,
			User:            &UserFormated{},
			Contents:        []*PostContentFormated{},
			CommentCount:    p.CommentCount,
			LatestRepliedOn: p.LatestRepliedOn,
			CreatedOn:       p.CreatedAt,
			ModifiedOn:      p.UpdatedAt,
			Tags:            tagsMap,
			IPLoc:           p.IPLoc,
		}
	}

	return nil
}

func (p *Post) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*Post, error) {
	var posts []*Post
	var err error
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if p.UserID > 0 {
		db = db.Where("user_id = ?", p.UserID)
	}
	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}

	if err = db.Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Post) Count(db *gorm.DB, conditions *ConditionsT) (int64, error) {
	var count int64
	if p.UserID > 0 {
		db = db.Where("user_id = ?", p.UserID)
	}
	for k, v := range *conditions {
		if k != "ORDER" {
			db = db.Where(k, v)
		}
	}
	if err := db.Model(p).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func FormatPosts(posts []*Post) ([]*PostFormated, error) {

	postIds := []int64{}
	userIds := []int64{}
	for _, post := range posts {
		postIds = append(postIds, int64(post.ID))
		userIds = append(userIds, post.UserID)
	}

	postContents, err := GetPostContentsByIDs(postIds)
	if err != nil {
		return nil, err
	}

	users, err := GetUsersByIDs(userIds)
	if err != nil {
		return nil, err
	}

	// 数据整合
	postsFormated := []*PostFormated{}
	for _, post := range posts {
		postFormated := post.Format()

		for _, user := range users {
			if user.ID == postFormated.UserID {
				postFormated.User = user.Format()
			}
		}
		for _, content := range postContents {
			if content.ID == post.ID {
				postFormated.Contents = append(postFormated.Contents, content.Format())
			}
		}

		postsFormated = append(postsFormated, postFormated)
	}

	return postsFormated, nil

}

func (p *PostContent) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*PostContent, error) {
	var contents []*PostContent
	var err error
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if p.PostID > 0 {
		db = db.Where("id = ?", p.PostID)
	}

	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}

	if err = db.Find(&contents).Error; err != nil {
		return nil, err
	}

	return contents, nil
}

/* 创建Post */
func CreatePostItem(post *Post) (uint, error) {
	post.LatestRepliedOn = time.Now().Unix()
	err := dao.DB.Debug().Create(&post).Error
	return post.ID, err
}

/* 创建Post Content */
func CreatePostContent(postContent *PostContent) (err error) {
	err = dao.DB.Debug().Create(&postContent).Error
	return
}

func GetPost(id int64) (*PostFormated, error) {
	post, err := GetPostByID(id)

	if err != nil {
		return nil, err
	}

	postContents, err := GetPostContentsByIDs([]int64{int64(post.ID)})
	if err != nil {
		return nil, err
	}

	users, err := GetUsersByIDs([]int64{post.UserID})
	if err != nil {
		return nil, err
	}

	// 数据整合
	postFormated := post.Format()
	for _, user := range users {
		postFormated.User = user.Format()
	}
	for _, content := range postContents {
		if content.PostID == int64(post.ID) {
			postFormated.Contents = append(postFormated.Contents, content.Format())
		}
	}
	return postFormated, nil
}

func GetPosts(conditions *ConditionsT, offset, limit int) (postList []*Post, err error) {
	return (&Post{}).List(dao.DB, conditions, offset, limit)
}

func GetPostByID(id int64) (post *Post, err error) {
	post = new(Post)
	if err := dao.DB.Debug().Where("id = ?", id).First(post).Error; err != nil {
		return nil, err
	}
	return
}

func GetPostContentsByIDs(ids []int64) ([]*PostContent, error) {
	return (&PostContent{}).List(dao.DB, &ConditionsT{
		"post_id IN ?": ids,
		"ORDER":        "sort ASC",
	}, 0, 0)
}

func GetPostCount(conditions *ConditionsT) (int64, error) {
	return (&Post{}).Count(dao.DB, conditions)
}

func GetPostList(req *PostListReq) ([]*PostFormated, error) {
	posts, err := GetPosts(req.Conditions, req.Offset, req.Limit)

	if err != nil {
		return nil, err
	}

	return FormatPosts(posts)
}

func DeletePostItem(id string) (err error) {
	err = dao.DB.Debug().Where("id=?", id).Delete(&Post{}).Error
	return
}
