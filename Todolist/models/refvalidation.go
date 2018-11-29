package models

import "wastest/common"

//해당 작업에 선택한 작업을 참조 시킬 수 있는지 확인
func Checkref(parent_id, child_id uint) (bool, error) {

	var parent Todo
	var child Todo
	var err error
	err = parent.FindById(parent_id)
	if err != nil {
		return false, err
	}
	err = child.FindById(child_id)
	if err != nil {
		return false, err
	}
	if parent.Done == "N" && child.Done == "Y" {
		return false, nil
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
