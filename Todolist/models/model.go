package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
	"wastest/common"
)

type Todo struct {
	Id        uint `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string `gorm:"type:varchar(1020)"`
	Done      string `gorm:"default:'N'"`

	Parents  []*Todo `gorm:"many2many:parents;association_jointable_foreignkey:todos_id"`
	Children []*Todo `gorm:"many2many:children;association_jointable_foreignkey:todos_id"`
}

func (todo *Todo) BeforeCreate(scope *gorm.Scope) error {

	scope.SetColumn("created_at", time.Now())

	scope.SetColumn("updated_at", time.Now())

	return nil

}
func (todo *Todo) BeforeUpdate(scope *gorm.Scope) error {

	//scope.SetColumn("updated_at", time.Now())

	return nil
}
func (todo *Todo) AfterDelete(scope *gorm.Scope) error {
	db := scope.DB()
	deleted_id := todo.Id
	//해당 작업이 삭제 될 경우 연결 돼 있는 정보를 모두 분리 해 준다.
	if err := db.Exec("delete from children where todos_id = ? or todo_id=?", deleted_id, deleted_id).Error; err != nil {
		return err
	}
	if err := db.Exec("delete from parents where todos_id = ? or todo_id=?", deleted_id, deleted_id).Error; err != nil {
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
	if f, err := todo.SameCountRefTodo(todo.Children); f == false {
		return err
	}
	if err := db.Create(&newtodo).Error; err != nil {
		return err
	}
	if f, _ := todo.IsPossibleConnect(); f == false {
		return errors.New("invalid value")

	}
	todo.Id = newtodo.Id
	fmt.Println(newtodo)
	fmt.Println(todo)
	todo.Connectref()

	return nil
}

func (todo *Todo) Insertreflist() error {

	db := common.GetDB()
	p := todo.Parents
	c := todo.Children
	/*
		if err := db.Model(&Todo{Id: todo.Id}).Association("Reflist").Clear().Error; err != nil {
			fmt.Println(err)
			return err
		}*/
	if err := db.Model(&Todo{Id: todo.Id}).Association("Parents").Append(&p).Error; err != nil {
		return err
	}
	if err := db.Model(&Todo{Id: todo.Id}).Association("Children").Append(&c).Error; err != nil {
		return err
	}
	return nil

}

func (todo *Todo) UpdateTodo() error {
	db := common.DB
	origin := Todo{
		Id: todo.Id,
	}
	db.Find(&origin, todo.Id)
	inter := intersection(todo.Children, origin.Children)
	intermap := TodoListToMap(inter)
	takoff := origin.Diffset(intermap)
	addon := origin.Diffset(intermap)
	//다시봐야함
	newtodo := Todo{

		Id: todo.Id,
	}
	for _, do := range inter {
		newtodo.Children = append(newtodo.Children, do)
	}
	for _, do := range addon {
		newtodo.Children = append(newtodo.Children, do)
	}

	if f := newtodo.ExistRef(takoff); f == false {
		return errors.New("invalid value")
	}
	if x, _ := newtodo.IsPossibleConnect(); x == false {
		return errors.New("invalid value")
	}

}
func TodoListToMap(todolist []*Todo) map[uint]bool {

	ret := make(map[uint]bool)
	for _, todo := range todolist {
		ret[todo.Id] = true
	}
	return ret
}
func (todo *Todo) Diffset(inter map[uint]bool) []*Todo {
	ret := make([]*Todo, 0)

	for _, child := range todo.Children {
		if _, ok := inter[child.Id]; ok {
			continue
		} else {
			ret = append(ret, child)
		}
	}
	return ret

}

//해당 작업이 완료 상태가 될수 있는지체크
func (todo *Todo) CheckDone() (bool, error) {
	descendant := todo.FindFamiliy("Children")
	for _, ref := range descendant {
		if ref.Done == "N" {
			return false, nil
		}
	}
	return true, nil

}

//각 작업끼리 연결 돼 있는지 확인 하는 메소드
func (todo *Todo) ExistRef(reflist []*Todo) bool {
	db := common.GetDB()

	childlist := make(map[uint]bool)
	db.Preload("Children").Find(&todo, "id in (?)", todo.Id)
	for _, ref := range todo.Children {
		childlist[ref.Id] = true
	}
	for _, ref := range reflist {
		if _, ok := childlist[ref.Id]; ok {
			continue
		} else {
			return false
		}

	}
	return true

}

//참조를 하기 위한 메소드
func (todo *Todo) Connectref() error {

	db := common.GetDB()

	for _, ref := range todo.Children {
		//todo.Children = append(todo.Children, ref)
		//ref.Parents = append(ref.Parents, todo)
		if err := db.Model(&Todo{Id: todo.Id}).Association("Children").Append(ref).Error; err != nil {
			return err
		}
		if err := db.Model(&Todo{Id: ref.Id}).Association("Parents").Append(todo).Error; err != nil {
			return err
		}
	}

	return nil

}

//연결된 작업을 제거 한느 메소드
func (todo *Todo) Disconnecref(reflist []*Todo) error {

	db := common.GetDB()

	if todo.ExistRef(todo.Children) == true {
		for _, ref := range reflist {

			if err := db.Model(&Todo{Id: todo.Id}).Association("Children").Delete(ref).Error; err != nil {
				return err
			}
			if err := db.Model(&Todo{Id: ref.Id}).Association("Parents").Delete(todo).Error; err != nil {
				return err
			}

		}

	}
	return nil

}

// 참조 하는 작업이 DB에 존재하는지 확인 하는 메소드
func (todo *Todo) SameCountRefTodo(reflist []*Todo) (res bool, err error) {
	db := common.GetDB()

	var qeurycount int
	refcount := len(reflist)

	if refcount > 0 {
		querylist := make([]uint, 0)
		for _, ref := range reflist {
			querylist = append(querylist, ref.Id)
		}

		db.Model(&Todo{}).Where("id in (?)", querylist).Count(&qeurycount)
		if refcount != qeurycount {
			return false, errors.New("Invaild ref")
		}
	} else {
		return true, nil
	}
	return true, nil

}

//참조가 가능한지 확인 하는 메소드
func (todo *Todo) IsPossibleConnect() (res bool, err error) {

	ancestor := todo.FindFamiliy("Parents")

	if res, err = todo.SameCountRefTodo(todo.Children); res == false {
		return
	}
	for _, child := range todo.Children {
		ancestor = append(ancestor, child)
	}

	res = true
	for _, ref := range todo.Children {
		descendant := ref.FindFamiliy("Children")
		inter := intersection(ancestor, descendant)
		if len(inter) > 0 {
			return false, nil
		}
		if todo.Done == "Y" {
			for _, child_ref := range descendant {
				if child_ref.Done == "N" {
					return false, nil

				}
			}
		}
	}
	return res, nil

}

//부모 또는 자식 노드 리스트를 찾아주는 메소드
func (todo *Todo) FindFamiliy(familytype string) (res []*Todo) {

	db := common.DB
	ret := make([]*Todo, 0)
	st := make([]uint, 0)
	familyids := make(map[uint]*Todo)
	st = append(st, todo.Id)
	fmt.Println(todo.Id)

	for {

		db.Preload(familytype).Find(&ret, "id in (?)", st)
		if len(ret) == 0 {
			break
		}

		st = make([]uint, 0)

		for _, curr := range ret {
			child := reflect.ValueOf(curr).Elem()
			next := child.FieldByName(familytype).Interface().([]*Todo)
			familyids[curr.Id] = curr

			for _, ref := range next {
				st = append(st, ref.Id)

			}

		}

		ret = make([]*Todo, 0)

	}
	st = make([]uint, 0)
	for _, value := range familyids {

		res = append(res, value)

	}
	return

}

//교집합을 구하는 함수
func intersection(a, b []*Todo) (inter []*Todo) {
	// interacting on the smallest list first can potentailly be faster...but not by much, worse case is the same
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	done := false
	for i, l := range low {
		for j, h := range high {
			// get future index values
			f1 := i + 1
			f2 := j + 1
			if l.Id == h.Id {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					// if the future values aren't the same then that's the end of the intersection
					if low[f1] != high[f2] {
						done = true
					}
				}
				// we don't want to interate on the entire list everytime, so remove the parts we already looped on will make it faster each pass
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		// nothing in the future so we are done
		if done {
			break
		}
	}
	return
}
