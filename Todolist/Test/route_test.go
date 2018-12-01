package Test

import (
	"fmt"
	"testing"
	"wastest/Todolist/Route"
)

func TestConvert(t *testing.T) {
	routtodo := Route.Todo{Id: 1, Children: []uint64{1, 2}}
	fmt.Println(routtodo.ConvertToModel())

}
