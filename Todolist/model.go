package Todolist

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testwas/common"
	"time"
)


type Todo struct{

	Id uint `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title string `gorm:"type:varchar(100);unique"`
	Done string `gorm:"default:'N'"`
	Reflist []*Todo `gorm:"many2many:ref;association_jointable_foreignkey:todos_id"`
}


func (todo *Todo) BeforeCreate(scope *gorm.Scope) (err error){
	scope.SetColumn("created_at",time.Now())
	scope.SetColumn("updated_at",time.Now())
	fmt.Println("run")
	return
}
func (todo *Todo) BeforeUpdate(scope *gorm.Scope) (err error){
	scope.SetColumn("updated_at",time.Now())

	return
}
func (todo *Todo) AfterDelete(scope *gorm.Scope)(err error){
	db := scope.DB()
	deleted_id := todo.Id
	//해당 작업이 삭제 될 경우 연결 돼 있는 정보를 모두 분리 해 준다.
	db.Exec("delete from ref where todos_id = ? or todo_id=?",deleted_id,deleted_id)
	return

}
func (todo *Todo)CreateTodo() error{
	db := common.GetDB()
	err := db.Create(&todo).Error
	if err == nil{
		return nil
	}else{
		return err
	}
}

func (todo *Todo)Addref(todolist []*Todo){
	db := common.GetDB()
	var update_todo Todo

	db.Find(&update_todo,"Id = ?", todo.Id)
	if todo.Done == "Y"{

	}else{

	}

}
func IsAbleDone(id int) (bool, error){
	var count int
	var todolist []Todo
	db := common.GetDB()
	err := db.Preload("Todolist","Done = ?","N").Find(&todolist,"id = ?",id).Count(&count).Error
	if err == nil{
		if count == 0{
			return true,nil
		}else{
			return false,nil
		}

	}else{
		return false,err
	}
}
