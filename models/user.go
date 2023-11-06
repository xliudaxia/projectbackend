package models

import (
	"bubble/dao"
	"bubble/pkg/errcode"
	"log"

	"gorm.io/gorm"
)

const (
	UserStatusNormal int = iota + 1
	UserStatusClosed
)

// 用户结构体
type User struct {
	ID       int64  `json:id`
	NickName string `json:"nickname"`
	UserName string `json:username`
	Phone    int    `json:phone`
	PassWord string `json:password`
	Salt     string `json:"salt"`
	Status   int    `json:"status"`
	Email    string `json:email`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
	Deleted  bool   `json:deleted`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type ChangePasswordReq struct {
	Password    string `json:"password" form:"password" binding:"required"`
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
}

type ChangeNicknameReq struct {
	Nickname string `json:"nickname" form:"nickname" binding:"required"`
}

type AuthRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserFormated struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickname"`
	UserName string `json:"username"`
	Status   int    `json:"status"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type UserInfo struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Phone    int    `json:"phone"`
	Email    string `json:"email"`
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

func NewCreateUser(user *User) (*User, error) {
	return user, dao.DB.Debug().Create(&user).Error
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

func (u *User) Format() *UserFormated {
	if u != nil {
		return &UserFormated{
			ID:       u.ID,
			UserName: u.UserName,
			NickName: u.NickName,
			Avatar:   u.Avatar,
			IsAdmin:  u.IsAdmin,
		}
	}

	return nil
}

func GetUserByUserName(username string) (user *User, err error) {
	user = new(User)
	if err = dao.DB.Debug().Where("user_name = ?", username).First(user).Error; err != nil {
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

func (u *User) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*User, error) {
	var users []*User
	var err error
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}

	if err = db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByUsername(username string) (*User, error) {
	user, err := GetUserByUserName(username)
	if err != nil {
		return nil, err
	}
	if user.ID > 0 {
		return user, nil
	}
	return nil, errcode.NoExistUsername
}

func GetUsersByIDs(ids []int64) ([]*User, error) {
	user := &User{}

	return user.List(dao.DB, &ConditionsT{
		"id IN ?": ids,
	}, 0, 0)
}
