package models

import "bubble/dao"

// Todo model
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

/*
Todo增删改查
*/
func CreateTodoItem(todo *Todo) (err error) {
	err = dao.DB.Debug().Create(&todo).Error
	return
}

func GetTodoList() (todolist []*Todo, err error) {
	if err = dao.DB.Debug().Find(&todolist).Error; err != nil {
		return nil, err
	}
	return
}

func GetTodoItem(id string) (todo *Todo, err error) {
	todo = new(Todo)
	if err = dao.DB.Debug().Where("id = ?", id).First(todo).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateTodoItem(todo *Todo) (err error) {
	err = dao.DB.Debug().Save(todo).Error
	return
}

func DeleteTodoItem(id string) (err error) {
	err = dao.DB.Debug().Where("id=?", id).Delete(&Todo{}).Error
	return
}
