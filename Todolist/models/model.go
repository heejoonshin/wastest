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
	//db := scope.DB()
	//deleted_id := todo.Id
	//해당 작업이 삭제 될 경우 연결 돼 있는 정보를 모두 분리 해 준다.
	//db.Exec("delete from ref where todos_id = ? or todo_id=?", deleted_id, deleted_id)
	return
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
func (todo *Todo) Insertreflist() {

	db := common.GetDB()
	reflist := todo.Reflist

	db.Model(&Todo{Id: todo.Id}).Association("Reflist").Append(&reflist)

}

/*
func (todo *Todo) Delref(childlist []*Todo) (bool,error){
	if err := todo.GetReflist(); err != nil{
		return false,err
	}
	DeleteMap := make(map[uint]bool)
	for _,child := range childlist{
		if err := child.FindById(); err != nil{
			return false,err
		}
		//DeleteMap
	}
}*/
/*
func (todo *Todo) Addref(childlist []*Todo) (bool,error){
	//db := common.GetDB()
	if err := todo.GetAllReflist(); err != nil{
		return false,err
	}
	for _,child := range childlist{
		if err := child.FindById(); err != nil{
			return false,err
		}
		if todo.Done == "N" && child.Done == "Y"{
			return false,errors.New("Invalid Ref")
		}
		check,err:=VaildIntersect(todo,child)
		if err != nil{
			return false,err
		}
		if !check{
			return false,nil
		}
		todo.Reflist = append(todo.Reflist,&Todo{Id:child.Id})
		todo.RefAll = append(todo.RefAll,&Todo{Id:child.Id})
		for _,childref := range child.Reflist{
			todo.RefAll = append(todo.RefAll,&Todo{Id:childref.Id})
		}
	}
	return true,nil
}*/

/*
func (todo *Todo)ValidRefs() (bool,error){
	db := common.GetDB()
	reflist := todo.Reflist

	db.Preload("Reflist").Find(&reflist)



}*/
