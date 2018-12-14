package models

import (
	"errors"
	"fmt"
	"github.com/heejoonshin/wastest/common"
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
)

type Todo struct {
	Id        uint64 `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string  `gorm:"type:varchar(1020)"`
	Done      string  `gorm:"default:'N'"`
	Parents   []*Todo `gorm:"many2many:parents;association_jointable_foreignkey:todos_id"`
	Children  []*Todo `gorm:"many2many:children;association_jointable_foreignkey:todos_id"`
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
	if err := db.Exec("delete from children where todos_id = ? or todo_id=?", deleted_id, deleted_id).Error; err != nil {
		return err
	}
	if err := db.Exec("delete from parents where todos_id = ? or todo_id=?", deleted_id, deleted_id).Error; err != nil {
		return err
	}
	return nil
}
func (todo *Todo) DelTodo() error {
	db := common.DB
	if err := db.Delete(&todo, "id = ?", todo.Id).Error; err != nil {
		return err

	}
	return nil

}

func (todo *Todo) FindById() error {
	db := common.GetDB()
	err := db.Preload("Children").Find(&todo, "id =?", todo.Id).Error
	if err != nil {
		return err
	}
	todo.FindAllInfo()

	return nil
}
func (todo *Todo) FindAllInfo() error {
	db := common.GetDB()

	idxs := make([]uint64, 0)

	for _, child := range todo.Children {
		idxs = append(idxs, child.Id)
	}

	if err := db.Find(&todo.Children, "id in (?)", idxs).Error; err != nil {
		return err
	}
	return nil

}

func (todo *Todo) CreateTodo() error {
	db := common.GetDB()
	newtodo := Todo{
		Title: todo.Title,
	}
	todo.Done = "N"
	//데이터 베이스에 참조 하려는 일의 아이디가 있는지 확인
	if err := todo.SameCountRefTodo(todo.Children); err != nil {
		return err
	}
	//중복로직 (구현 실수)
	if err := todo.FindAllInfo(); err != nil {
		return err
	}

	//참조를 걸수 있는지 확인 하는 로직
	if err := todo.IsPossibleConnect(); err != nil {
		return err

	}
	//완료를 할 수 있는 로직
	if err := todo.ispossibleDone(); err != nil {
		return err
	}
	//실제로 데이터 베이스에 데이터를 생성 하는 로직
	if err := db.Create(&newtodo).Error; err != nil {
		return err
	}

	// 데이터를 디비에 생성하면 아이디를 얻을 수 있다.
	todo.Id = newtodo.Id
	fmt.Println(newtodo)
	fmt.Println(todo)
	//참조된 일을 연결
	todo.Connectref()

	db.Preload("Children").Find(&todo, "id =?", todo.Id)
	todo.FindAllInfo()

	return nil
}

func (todo *Todo) Insertreflist() error {

	db := common.GetDB()

	c := todo.Children
	/*
		if err := db.Model(&Todo{Id: todo.Id}).Association("Reflist").Clear().Error; err != nil {
			fmt.Println(err)
			return err
		}*/

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
	if err := db.Preload("Parents").Find(&origin, todo.Id).Error; err != nil {
		return err
	}
	if err := db.Preload("Children").Find(&origin, todo.Id).Error; err != nil {
		return err
	}
	//바뀔 데이터중 삭제할 참조 데이터와 남겨둘 참조 데이터를 구하는 로직
	//바꾸려는 데이터에서 교집합을 빼면 추가될 데이터, 기존의 있는 데이터에서 교집합을 빼면 삭제될 데이터
	inter := intersection(todo.Children, origin.Children)
	intermap := TodoListToMap(inter)
	//db에 있는 정보에서 교집합을 빼는 부분
	takoff, err := origin.Diffset(intermap)
	if err != nil {
		return err
	}

	//바뀔 정보에서 교집합을 빼서 새롭게 참조되는 참조를 추가 하는 부분
	addon, err := todo.Diffset(intermap)
	if err != nil {
		return err
	}

	newtodo := Todo{

		Id:    todo.Id,
		Title: todo.Title,
		Done:  todo.Done,
	}

	for _, do := range inter {
		newtodo.Children = append(newtodo.Children, do)
	}
	for _, do := range addon {
		newtodo.Children = append(newtodo.Children, do)
	}
	newtodo.Parents = origin.Parents

	//디비에 있는 정보로 검색
	if f := origin.ExistRef(takoff); f == false {
		return errors.New("?")
	}
	//새롭게 업데이트 될 정보로 검색
	if err := newtodo.IsPossibleConnect(); err != nil {
		return err
	}
	if err := newtodo.ispossibleDone(); err != nil {
		return err
	}
	if err := db.Model(&Todo{}).Where("id =?", newtodo.Id).Update(&Todo{Title: newtodo.Title, Done: newtodo.Done}).Error; err != nil {
		return err
	}

	origin.Disconnecref(takoff)
	newtodo.Connectref()
	db.Preload("Children").Find(&todo, "id =?", newtodo.Id)
	newtodo.FindAllInfo()
	return nil

}
func TodoListToMap(todolist []*Todo) map[uint64]bool {

	ret := make(map[uint64]bool)
	for _, todo := range todolist {
		ret[todo.Id] = true
	}
	return ret
}
func (todo *Todo) Diffset(inter map[uint64]bool) ([]*Todo, error) {
	ret := make([]*Todo, 0)

	for _, child := range todo.Children {

		if _, ok := inter[child.Id]; ok {
			continue
		} else {
			if err := child.FindById(); err != nil {
				return ret, err
			}
			ret = append(ret, child)
		}
	}
	return ret, nil

}

//각 작업끼리 연결 돼 있는지 확인 하는 메소드
func (todo *Todo) ExistRef(reflist []*Todo) bool {
	db := common.GetDB()
	var temptodo Todo

	childlist := make(map[uint64]bool)
	db.Preload("Children").Find(&temptodo, "id in (?)", todo.Id)
	for _, ref := range temptodo.Children {
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
		if err := db.Model(&Todo{Id: ref.Id}).Association("Parents").Append(&Todo{Id: todo.Id}).Error; err != nil {
			return err
		}

	}

	return nil

}

//연결된 작업을 제거 하는 메소드
func (todo *Todo) Disconnecref(reflist []*Todo) error {

	db := common.GetDB()

	if todo.ExistRef(todo.Children) == true {
		for _, ref := range reflist {

			if err := db.Model(&Todo{Id: todo.Id}).Association("Children").Delete(ref).Error; err != nil {
				return err
			}
			if err := db.Model(&Todo{Id: ref.Id}).Association("Parents").Delete(&Todo{Id: todo.Id}).Error; err != nil {
				return err
			}

		}

	}
	return nil

}

// 참조 하는 작업이 DB에 존재하는지 확인 하는 메소드
func (todo *Todo) SameCountRefTodo(reflist []*Todo) (err error) {
	db := common.GetDB()

	var qeurycount int
	refcount := len(reflist)

	if refcount > 0 {
		querylist := make([]uint64, 0)
		for _, ref := range reflist {
			querylist = append(querylist, ref.Id)
		}

		db.Model(&Todo{}).Where("id in (?)", querylist).Count(&qeurycount)
		if refcount != qeurycount {
			return errors.New("선택한 일정이 존재 하지 않습니다.")
		}
	} else {
		return nil
	}
	return nil

}

//완료여부를 수정할 수 있는지
func (todo *Todo) ispossibleDone() error {
	//현재 작업을 N-> Y로 바꿀 때 부모의 작업이 N인 경우가 있으면 안된다.
	for _, p := range todo.Parents {
		if todo.Done == "Y" && p.Done == "N" {
			return errors.New("참조한 완료된 후에 완료할 수 있습니다.")
		}
	}
	//현재 작업을 Y -> N로 바꿀때 자식이 Y인 경우가 있으면 안되다.
	for _, c := range todo.Children {
		if todo.Done == "N" && c.Done == "Y" {
			return errors.New("잘못된 참조 영역이 있습니다.")

		}
	}
	return nil

}

//참조가 가능한지 확인 하는 메소드
func (todo *Todo) IsPossibleConnect() error {

	//참조할 데이터가 존재하는지 확인 하는 로직
	if err := todo.SameCountRefTodo(todo.Children); err != nil {
		return err
	}

	for _, ref := range todo.Children {
		descendant := ref.FindFamiliy("Children")

		for _, child_ref := range descendant {

			if todo.Id == child_ref.Id {
				return errors.New("사이클이 존재합니다. 영원히 일을 끝낼 수 없습니다.")
			}
		}
	}
	return nil

}

//부모 또는 자식 노드 리스트를 찾아주는 메소드
func (todo *Todo) FindFamiliy(familytype string) (res []*Todo) {

	db := common.DB
	ret := make([]*Todo, 0)
	st := make([]uint64, 0)
	familyids := make(map[uint64]*Todo)
	st = append(st, todo.Id)
	fmt.Println(todo.Id)

	for {

		db.Preload(familytype).Find(&ret, "id in (?)", st)
		if len(ret) == 0 {
			break
		}

		st = make([]uint64, 0)

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
	st = make([]uint64, 0)
	for _, value := range familyids {

		res = append(res, value)

	}
	return

}

//교집합을 구하는 함수
func intersection(a, b []*Todo) (inter []*Todo) {
	m := make(map[uint64]bool)

	for _, item := range a {
		m[item.Id] = true
	}

	for _, item := range b {
		if _, ok := m[item.Id]; ok {
			inter = append(inter, item)
		}
	}

	return
}
