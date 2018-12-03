package Test

import (
	"fmt"
	"testing"
	"time"
	"wastest/Todolist/Route"
	"wastest/Todolist/models"
)

func TestConvert(t *testing.T) {
	routtodo := Route.Todo{Id: 1, Children: []uint64{1, 2}}
	fmt.Println(routtodo.ConvertToModel())

}
func TestModelToView(t *testing.T) {
	A := &models.Todo{Id: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(), Done: "Y", Title: "test", Children: []*models.Todo{{Id: 3}}}
	fmt.Println(A)

	fmt.Println(Route.ModelToView(A))

}
