package models

import (
	"bubble/dao"

	"gorm.io/gorm"
)

// Phone model
type Phone struct {
	gorm.Model
	Name    string `json:"name" gorm:"not null"`
	Sex     string `json:"sex"`
	Phone   int64  `json:"phone" gorm:"not null"`
	WeChat  string `json:"wechat"`
	Label   string `json:"label"`
	QQNum   int64  `json:"qqnum"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Extra   string `json:"extra"`
	Status  bool   `json:"status"`
}

/*
Phone增删改查
*/
func CreatePhoneItem(phone *Phone) (err error) {
	err = dao.DB.Debug().Create(&phone).Error
	return
}

func GetPhoneList() (phonelist []*Phone, err error) {
	if err = dao.DB.Debug().Find(&phonelist).Error; err != nil {
		return nil, err
	}
	return
}

func GetPhoneItem(id string) (phone *Phone, err error) {
	phone = new(Phone)
	if err = dao.DB.Debug().Where("id = ?", id).First(phone).Error; err != nil {
		return nil, err
	}
	return
}

func UpdatePhoneItem(phone *Phone) (err error) {
	err = dao.DB.Debug().Save(phone).Error
	return
}

func DeletePhoneItem(id string) (err error) {
	err = dao.DB.Debug().Where("id=?", id).Delete(&Phone{}).Error
	return
}
