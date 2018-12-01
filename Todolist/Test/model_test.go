package Test

import (
	//"fmt"

	"fmt"
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
		//title := "test" + strconv.Itoa(i)
		//todo := models.Todo{Title: title}
		//todo.CreateTodo()

	}

}
func TestCreate(t *testing.T) {

	Init()
	//db:=common.GetDB()
	testcase := models.Todo{
		Title: "sssssvvs",
		Children: []*models.Todo{
			{
				Id: 7,
			},
		},
	}
	fmt.Println(testcase.IsPossibleConnect(testcase.Children))

}
