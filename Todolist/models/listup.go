package models

import (
	"wastest/common"
)

//페이징 처리
type Pageination struct {
	Limit int    `json:limit`
	Page  int    `json:page`
	Order string `json:order`
}
type result struct {
	Todolist  []*Todo
	Totaldata int
	Limit     int
	Page      int
}

func (pageination *Pageination) Listup() (res result, err error) {

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

	err = db.Preload("Children").Offset(offset).Limit(pageination.Limit).Order(pageination.Order).Find(&todolist).Error
	if err != nil {
		return
	}
	var totaldata int
	var temp []*Todo
	db.Find(&temp).Count(&totaldata)
	//totalpage /=pageination.Limit
	if err != nil {
		return res, err
	} else {
		for _, todo := range todolist {
			todo.FindAllInfo()
		}
		res = result{Todolist: todolist, Totaldata: totaldata, Limit: pageination.Limit, Page: pageination.Page}
		return res, nil
	}

}
