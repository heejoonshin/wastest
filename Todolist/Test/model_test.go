package Test

import (
	"strconv"
	"testing"
	"wastest/Todolist/models"
	"wastest/common"
)

func Init() {
	db := common.TestDBInit()
	db.AutoMigrate(&models.Todo{})
}
func InsertDumyData() {
	db := common.GetDB()
	tx := db.Begin()
	for i := 0; i < 10000; i++ {
		title := "test" + strconv.Itoa(i)
		tx.Save(&models.Todo{Title: title})
	}
	tx.Commit()
}
func TestXx(t *testing.T) {
	Init()
	InsertDumyData()

}
