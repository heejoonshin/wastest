package models

import (
	"fmt"
	"wastest/common"
)

/*
//해당 작업에 선택한 작업을 참조 시킬 수 있는지 확인
func VaildIntersect(parent, child *Todo) (bool, error) {


	var err error

	if parent.Done == "N" && child.Done == "Y" {
		return false, nil
	}
	err = child.GetAllReflist()
	if err != nil{
		return false,err
	}
	set := make(map[uint]bool)
	set[parent.Id] = true
	for _,item := range parent.Reflist{
		set[item.Id] =true
	}
	count := 0
	if _,ok := set[child.Id]; ok{
		return false,nil
	}
	for _, itme := range child.Reflist{
		if _,ok := set[itme.Id]; ok{
			count ++
			break
		}
	}
	if count > 0{
		return false,nil
	}
	return true, nil

}*/

func ValidationRef(U, V Todo) (bool, error) {

	queue := common.New()

	queue.PushBack(U)
	queue.PushBack(V)

	check := make(map[uint]bool)
	check[U.Id] = true
	check[V.Id] = true

	for queue.Len() > 0 {
		value := queue.Front()
		queue.PopFront()
		s, ok := value.(Todo)

		s.GetReflist()
		for _, child := range s.Reflist {
			if _, ok := check[child.Id]; ok {

			}

		}

		fmt.Println(s, ok)

	}
	return true, nil

}

//해당 작업이 완료 상태가 될수 있는지체크
func CheckDone(id int) (bool, error) {
	var count int
	var todolist []Todo
	db := common.GetDB()
	err := db.Preload("Reflist", "Done = ?", "N").Find(&todolist, "id = ?", id).Count(&count).Error
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
