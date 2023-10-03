package models

import (
	"bubble/dao"

	"gorm.io/gorm"
)

// Category model
type Category struct {
	gorm.Model
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Rank        int    `json:"rank"`
}

/*
category增删改查
*/
func CreateCategoryItem(category *Category) (err error) {
	err = dao.DB.Debug().Create(&category).Error
	return
}

func GetCategoryList() (categoryList []*Category, err error) {
	if err = dao.DB.Debug().Find(&categoryList).Error; err != nil {
		return nil, err
	}
	return
}

func GetCategoryItem(id string) (category *Category, err error) {
	category = new(Category)
	if err = dao.DB.Debug().Where("id = ?", id).First(category).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateCategoryItem(category *Category) (err error) {
	err = dao.DB.Debug().Save(category).Error
	return
}

func DeleteCategoryItem(id string) (err error) {
	err = dao.DB.Debug().Where("id=?", id).Delete(&Category{}).Error
	return
}
