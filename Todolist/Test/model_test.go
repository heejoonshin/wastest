package Test

import (
	//"fmt"

	"fmt"
	"strconv"
	"testing"
	"wastest/Todolist/models"
	"wastest/common"
)

func Init() {
	db := common.TestDBInit()
	//db.DropTableIfExists(&models.Todo{})
	db.AutoMigrate(&models.Todo{})
	//db.Model(&models.Todo{}).Delete(&models.Todo{})
}
func InsertDumyData(n int) {

	for i := 0; i < n; i++ {
		title := "test" + strconv.Itoa(i)
		todo := models.Todo{Title: title}
		todo.CreateTodo()

	}

}
func TestCreate2(t *testing.T) {
	Init()
	f := models.Todo{Title: "testing"}
	f.CreateTodo()
}
func TestCreate(t *testing.T) {

	Init()

	test := &models.Todo{
		Id:   14,
		Done: "Y",
		Children: []*models.Todo{
			{Id: 1},
		},
	}
	err := test.UpdateTodo()
	fmt.Println(err)

}

func TestDel(t *testing.T) {
	Init()

	test := &models.Todo{
		Id: 6,
	}
	test.DelTodo()
}
