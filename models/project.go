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
Project增删改查
*/
func CreateProjectItem(project *Project) (err error) {
	err = dao.DB.Debug().Create(&project).Error
	return
}

func GetProjectList() (projectlist []*Project, err error) {
	if err = dao.DB.Debug().Find(&projectlist).Error; err != nil {
		return nil, err
	}
	return
}

func GetProjectItem(id string) (project *Project, err error) {
	project = new(Project)
	if err = dao.DB.Debug().Where("id = ?", id).First(project).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateProjectItem(project *Project) (err error) {
	err = dao.DB.Debug().Save(project).Error
	return
}

func DeleteProjectItem(id string) (err error) {
	err = dao.DB.Debug().Where("id=?", id).Delete(&Project{}).Error
	return
}
