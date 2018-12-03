package Test

import (
	"errors"
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

//복수개의 참조가 있을 경우 테스트
func TestCreateChildren(t *testing.T) {

	Init()
	data1 := models.Todo{Title: "test1"}
	data2 := models.Todo{Title: "test2"}
	if err := data1.CreateTodo(); err != nil {
		t.Error("create modle err")
		//err
	}
	if err := data2.CreateTodo(); err != nil {
		t.Error("create modle err")

	}
	if data1.Title != "test1" {
		//err
		t.Error("create modle err")
	}
	if data2.Title != "test2" {
		//err
		t.Error("create modle err")
	}

	data3 := models.Todo{Title: "test3", Children: []*models.Todo{{Id: data1.Id}, {Id: data2.Id}}}
	if err := data3.CreateTodo(); err != nil {
		t.Error("참조가 있는 작업 생성 오류")
	}
	fmt.Println(data3)

}
func TestFailModifyDone(t *testing.T) {

	Init()
	data1 := models.Todo{Title: "test1"}
	data2 := models.Todo{Title: "test2"}
	if err := data1.CreateTodo(); err != nil {
		t.Error("작업 생성 실패")
		//err
	}
	if err := data2.CreateTodo(); err != nil {
		t.Error("작업 생성 실패")

	}
	if data1.Title != "test1" {
		//err
		t.Error("작업 생성 실패")
	}
	if data2.Title != "test2" {
		//err
		t.Error("작업 생성 실패")
	}

	data3 := models.Todo{Title: "test3", Children: []*models.Todo{{Id: data1.Id}, {Id: data2.Id}}}
	if err := data3.CreateTodo(); err != nil {
		t.Error("참조가 있는 작업 생성 오류")
	}

	data3.Done = "Y"
	if err := data3.UpdateTodo(); err == nil {
		t.Error("데이터 수정시 완료 하는 로직 에러: 참조된 작이 끝나지 않은 상황에서 완료로 변경")

	}

}
func testModletitle(todo *models.Todo, expact string) error {
	if err := todo.CreateTodo(); err != nil {
		return err
		//err
	}
	if todo.Title != expact {
		//err
		errors.New("작업 생성 실패")
	}
	return nil

}
func testModifyModle(todo *models.Todo, expact *models.Todo) error {

	if err := todo.UpdateTodo(); err != nil {
		return err
	}
	if len(todo.Children) != len(expact.Children) {
		return errors.New("업데이트 오류")
	}
	for i := 0; i < len(todo.Children); i++ {
		if todo.Children[i].Id != expact.Children[i].Id {
			return errors.New("업데이트 오류")
		}
	}
	return nil

}

//참조의 사이클이 생기는 데이터 테스트
func TestFailModifyCycle(t *testing.T) {

	Init()
	/*test Case
	데이터 1 -> 2 -> 3 -> 4 -> 5 연결 시키고 3->1번을 연결시켜 사이클을 만듬
	*/
	tc := make([]*models.Todo, 0)
	for i := 1; i <= 5; i++ {
		title := "test" + strconv.Itoa(i)
		GenModle := &models.Todo{Title: title}
		tc = append(tc, GenModle)
		if err := testModletitle(GenModle, title); err != nil {
			t.Error(err)
		}
	}
	for i := 1; i < 5; i++ {
		tc[i].Children = append(tc[i].Children, &models.Todo{Id: tc[i-1].Id})
		if err := testModifyModle(tc[i], &models.Todo{Children: []*models.Todo{{Id: tc[i-1].Id}}}); err != nil {
			t.Error(err)
		}
	}
	tc[2].Children = append(tc[2].Children, &models.Todo{Id: tc[0].Id})
	if err := testModifyModle(tc[2], &models.Todo{Children: []*models.Todo{{Id: tc[0].Id}, {Id: tc[3].Id}}}); err == nil {
		t.Error("사이클인 존재 하는 데이터가 통과")
	}

}

func TestDel(t *testing.T) {
	Init()
	tc := models.Todo{Title: "test1"}

	if err := testModletitle(&tc, "test1"); err != nil {
		t.Error(err)
	}
	id := tc.Id
	if err := tc.DelTodo(); err != nil {
		t.Error(err)
	}
	check := models.Todo{Id: id}
	if err := check.FindById(); err == nil {
		t.Error("삭제 에러")
	}
}
