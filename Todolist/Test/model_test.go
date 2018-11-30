package Test

import (
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
func TestCreate(t *testing.T) {

	Init()
	//InsertDumyData(10)

	db := common.DB

	//Case 1
	testTodo := models.Todo{Id: 1}

	if err := db.Find(&testTodo).Error; err != nil {

		t.Error(err)
	} else {
		fmt.Println(testTodo)

		if testTodo.Title != "test0" {
			t.Fail()

		}
	}
	//Case 2
	testTodo = models.Todo{
		Title: "reftest",
		Reflist: []*models.Todo{
			{
				Id: 1,
			},
			{
				Id: 2,
			},
		},
	}
	testTodo.CreateTodo()

}
