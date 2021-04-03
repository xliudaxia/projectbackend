package models

import (
	"bubble/dao"
	"time"
)

// Project model
type Project struct {
	ID                 int    `json:"id" xorm:"not null pk autoincr INT(20)"`
	ProjectName        string `json:"projectname"`
	ProjectDescription string `json:"projectdescription"`
	ProjectAdmin       string `json:"projectadmin"`
	ProjectStatus      string `json:"projectstatus"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

/*
CreateProjectItem 增加project记录
*/
func CreateProjectItem(project *Project) (err error) {
	err = dao.DB.Debug().Create(&project).Error
	return
}

/*
DeleteProjectItem 删除project记录
*/
func DeleteProjectItem(id string) (err error) {
	err = dao.DB.Debug().Where("id=?", id).Delete(&Project{}).Error
	return
}

/*
UpdateProjectItem 修改project记录
*/
func UpdateProjectItem(project *Project) (err error) {
	err = dao.DB.Debug().Save(project).Error
	return
}

/*
GetProjectList 查询project记录列表（返回数组）
*/
func GetProjectList() (projectlist []*Project, err error) {
	if err = dao.DB.Debug().Find(&projectlist).Error; err != nil {
		return nil, err
	}
	return
}

/*
GetProjectItem 查询单条project记录
*/
func GetProjectItem(id string) (project *Project, err error) {
	project = new(Project)
	if err = dao.DB.Debug().Where("id = ?", id).First(project).Error; err != nil {
		return nil, err
	}
	return
}

/*
QueryProjectByRule 根据条件查询project记录
(目前支持根据项目名及项目介绍信息进行检索)
*/
func QueryProjectByRule(keyword string) (projectlist []*Project, err error) {
	// project = new(Project)
	if keyword == "" {
		if err = dao.DB.Debug().Find(&projectlist).Error; err != nil {
			return nil, err
		}
		return
	} else {
		if err = dao.DB.Debug().Where("project_name LIKE ?", "%"+keyword+"%").Or("project_description LIKE ?", "%"+keyword+"%").Find(&projectlist).Error; err != nil {
			return nil, err
		}
		return
	}
}
