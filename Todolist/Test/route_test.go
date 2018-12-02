package Test

import (
	"fmt"
	"testing"
	"wastest/Todolist/Route"
	"wastest/Todolist/models"
)

func TestConvert(t *testing.T) {
	routtodo := Route.Todo{Id: 1, Children: []uint64{1, 2}}
	fmt.Println(routtodo.ConvertToModel())

}
func TestModelToView(t *testing.T) {
	var A []*models.Todo
	A = append(A, &models.Todo{Id: 1})
	Route.ModelToView(&A)
	B := models.Todo{Id: 1}
	Route.ModelToView(&B)

}
