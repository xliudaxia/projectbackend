package models

import (
	"bubble/dao"

	"gorm.io/gorm"
)

// Link model
type Link struct {
	gorm.Model
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Rank        int    `json:"rank"`
	Public      int    `json:"public"`
	Status      int    `json:"status"`
	CID         int64  `json:"cid"`
}

/*
Link增删改查
*/
func CreateLinkItem(link *Link) (err error) {
	err = dao.DB.Debug().Create(&link).Error
	return
}

func GetLinkList() (linkList []*Link, err error) {
	if err = dao.DB.Debug().Find(&linkList).Error; err != nil {
		return nil, err
	}
	return
}

/** 查询指定用户链接 */
func GetUserLinkList(userId string) (linkList []*Link, err error) {
	if err = dao.DB.Debug().Where("user_id = ?", userId).Find(&linkList).Error; err != nil {
		return nil, err
	}
	return
}

func GetLinkItem(id string) (link *Link, err error) {
	link = new(Link)
	if err = dao.DB.Debug().Where("id = ?", id).First(link).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateLinkItem(link *Link) (err error) {
	err = dao.DB.Debug().Save(link).Error
	return
}

func DeleteLinkItem(id string) (err error) {
	err = dao.DB.Debug().Where("id=?", id).Delete(&Link{}).Error
	return
}
