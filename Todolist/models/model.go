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

func (todo *Todo) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("created_at", time.Now())
	scope.SetColumn("updated_at", time.Now())
	fmt.Println("run")
	return
}
func (todo *Todo) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("updated_at", time.Now())

	return
}
func (todo *Todo) AfterDelete(scope *gorm.Scope) (err error) {
	db := scope.DB()
	deleted_id := todo.Id
	//해당 작업이 삭제 될 경우 연결 돼 있는 정보를 모두 분리 해 준다.
	db.Exec("delete from ref where todos_id = ? or todo_id=?", deleted_id, deleted_id)
	return

}
func (todo *Todo) CreateTodo() error {
	db := common.GetDB()
	var validCount int
	reflist_count := len(todo.Reflist)

	var err error

	db.Find(&todo.Reflist).Count(validCount)
	if reflist_count != validCount {
		return errors.New("Invalid Count")

	}
	err = db.Create(&todo).Error
	if err == nil {
		return nil
	} else {
		return err
	}
}
func (todo *Todo) FindById(id uint) error {
	db := common.GetDB()
	err := db.Preload("Reflist").Find(&todo, id).Error
	if err == nil {
		return nil
	} else {
		return err
	}
}
func (todo *Todo) Addref(todolist []*Todo) {
	db := common.GetDB()
	var update_todo Todo

	db.Find(&update_todo, "Id = ?", todo.Id)
	if todo.Done == "Y" {

	} else {

	}

}

/*
func (todo *Todo)ValidRefs() (bool,error){
	db := common.GetDB()
	reflist := todo.Reflist

	db.Preload("Reflist").Find(&reflist)



}*/
