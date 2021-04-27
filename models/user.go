package models

import (
	"bubble/dao"
	"log"
)

// 用户结构体
type User struct {
	ID        int    `json:id`
	UserName  string `json:username`
	PassWord  string `json:password`
	Phone     int    `json:phone`
	WorkCode  string `json:workcode`
	Email     string `json:email`
	CompanyID string `json:companyid`
	Validated bool   `json:validated`
	Deleted   bool   `json:deleted`
}

type JwtToken struct {
	Token string `json:"token"`
}

type UserInfo struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Phone     int    `json:"phone"`
	WorkCode  string `json:"work_code"`
	CompanyID string `json:"company_id"`
	Email     string `json:"email"`
}

type UserRole struct {
	UserId  int    `json:"userId" xorm:"not null pk INT(11)"`
	RoleKey string `json:"roleKey" xorm:"not null pk VARCHAR(20)"`
}

//创建一个用户
func CreateUserItem(user *User) (err error) {
	err = dao.DB.Debug().Create(&user).Error
	return
}

//获取用户列表
func GetUserList() (userlist []*User, err error) {
	if err = dao.DB.Debug().Find(&userlist).Error; err != nil {
		return nil, err
	}
	return
}

//获取单个用户记录
func GetUserItem(id string) (user *User, err error) {
	user = new(User)
	if err = dao.DB.Debug().Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return
}

//更新用户信息
func UpdateUserItem(user *User) (err error) {
	err = dao.DB.Debug().Save(user).Error
	return
}

//删除单个用户
func DeleteUserItem(id string) (err error) {
	err = dao.DB.Debug().Where("id=?", id).Delete(&User{}).Error
	return
}

//QueryLoginUser: 获取当前登录用户信息
func QueryLoginUser(username string, action string) (user *User, err error) {
	user = new(User)
	err = dao.DB.Debug().Where(action+"=?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return
}

//注册用户信息校验
func RegisterUserCheck(username string, action string) (user *User, err error) {
	user = new(User)
	err = dao.DB.Debug().Where(action+"=?", username).First(user).Error
	if err != nil {
		log.Println("Find user Failed!", err.Error())
	}
	return user, nil
}
