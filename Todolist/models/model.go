package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	"time"
	"wastest/common"
)

type Todo struct {
	Id        uint `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string  `gorm:"type:varchar(100);unique"`
	Done      string  `gorm:"default:'N'"`
	Reflist   []*Todo `gorm:"many2many:ref;association_jointable_foreignkey:todos_id"`
}

func (todo *Todo) BeforeCreate(scope *gorm.Scope) error {

	scope.SetColumn("created_at", time.Now())

	scope.SetColumn("updated_at", time.Now())

	return nil

}
func (todo *Todo) BeforeUpdate(scope *gorm.Scope) error {

	scope.SetColumn("updated_at", time.Now())

	return nil
}
func (todo *Todo) AfterDelete(scope *gorm.Scope) error {
	db := scope.DB()
	deleted_id := todo.Id
	//해당 작업이 삭제 될 경우 연결 돼 있는 정보를 모두 분리 해 준다.
	if err := db.Exec("delete from ref where todos_id = ? or todo_id=?", deleted_id, deleted_id).Error; err != nil {
		return err
	}
	return nil
}

//해당 todo에 연결된 모든 자식 노드의 정보를 받아오는 함수
func (todo *Todo) GetReflist() error {
	db := common.GetDB()

	err := db.Preload("Reflist").Find(&todo).Error
	if err != nil {
		return err
	}
	return nil
}

func (todo *Todo) FindById() error {
	db := common.GetDB()
	err := db.Find(&todo).Error
	if err != nil {
		return err
	}
	return nil
}
func (todo *Todo) CreateTodo() error {
	db := common.GetDB()
	newtodo := Todo{
		Title: todo.Title,
	}
	refcount := len(todo.Reflist)
	var qeurycount int
	if len(todo.Reflist) > 0 {
		querylist := make([]uint, 0)
		for _, ref := range todo.Reflist {
			querylist = append(querylist, ref.Id)
		}

		db.Model(&Todo{}).Where("id in (?)", querylist).Count(&qeurycount)
		if refcount != qeurycount {
			return errors.New("Invaild ref")
		}
	}
	fmt.Println(newtodo)

	if err := db.Create(&newtodo).Error; err != nil {
		return err
	}
	if len(todo.Reflist) > 0 {
		if err := todo.Insertreflist(); err != nil {
			return err

		}
	}
	return nil

}
func (todo *Todo) Insertreflist() error {

	db := common.GetDB()
	reflist := todo.Reflist
	if err := db.Model(&Todo{Id: todo.Id}).Association("Reflist").Clear().Error; err != nil {
		fmt.Println(err)
		return err
	}
	if err := db.Model(&Todo{Id: todo.Id}).Association("Reflist").Append(&reflist).Error; err != nil {
		return err
	}
	return nil

}

func (todo *Todo) UpdateTodo() error {
	db := common.DB
	var updatedTodo Todo
	if err := db.Model(&updatedTodo).Where("id = ?", todo.Id).Update(todo).Error; err != nil {
		return err
	}
	return nil

}

//해당 작업이 완료 상태가 될수 있는지체크
func (todo *Todo) CheckDone() (bool, error) {
	var count int
	var todolist []Todo
	db := common.GetDB()
	err := db.Preload("Reflist", "Done = ?", "N").Find(&todolist).Count(&count).Error
	if err == nil {
		if count == 0 {
			return true, nil
		} else {
			return false, nil
		}

	} else {
		return false, err
	}
}
