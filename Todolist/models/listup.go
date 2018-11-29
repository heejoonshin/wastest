package models

import "wastest/common"

//페이징 처리
type Pageination struct {
	Limit int
	Page  int
	Order string
}

func (pageination *Pageination) Listup() ([]*Todo, error) {

	db := common.GetDB()
	offset := pageination.Limit * (pageination.Page - 1)
	var todolist []*Todo

	switch pageination.Order {
	case "recent":
		pageination.Order = "updated_at"
	case "done":
		pageination.Order = "done desc"
	default:
		pageination.Order = "id"

	}

	err := db.Preload("Reflist").Offset(offset).Limit(pageination.Limit).Order(pageination.Order).Find(&todolist).Error
	if err != nil {
		return todolist, err
	} else {
		return todolist, nil
	}

}
