package Route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heejoonshin/wastest/Todolist/models"
	"net/http"
	"strconv"
	"strings"
)

type Todo struct {
	Id       uint64   `json:id`
	Title    string   `json:title`
	Done     string   `json:done`
	Children []uint64 `json:children"`
}
type ViewStruct struct {
	Id        string
	CreatedAt string
	UpdatedAt string
	Title     string
	Done      string
}
type TodoListStruct struct {
	Todolist  []*ViewStruct
	Totaldata int
	Limit     int
	Page      int
}

func TodoGroup(router *gin.RouterGroup) {
	router.POST("/", CreateTodo)
	router.GET("/:id", GetTodo)
	router.PUT("/:id", ModfiyTodo)
	router.DELETE("/:id", DeleteTodo)

}

func TodolistGroup(router *gin.RouterGroup) {

	router.GET("/", GetTodolist)

}

//프론트에서 보여질 값으로 변환 하는 함수
func ModelToView(todo *models.Todo) (ret *ViewStruct) {
	childlist := ""
	ret = &ViewStruct{
		Id:        fmt.Sprint(todo.Id),
		CreatedAt: fmt.Sprintf("%d-%d-%d %d:%d:%d", todo.CreatedAt.Year(), todo.CreatedAt.Month(), todo.CreatedAt.Day(), todo.CreatedAt.Hour(), todo.CreatedAt.Minute(), todo.CreatedAt.Second()),
		UpdatedAt: fmt.Sprintf("%d-%d-%d %d:%d:%d", todo.UpdatedAt.Year(), todo.UpdatedAt.Month(), todo.UpdatedAt.Day(), todo.UpdatedAt.Hour(), todo.UpdatedAt.Minute(), todo.UpdatedAt.Second()),
		Title:     todo.Title,
		Done:      todo.Done,
	}
	for _, child := range todo.Children {
		childlist = childlist + "@" + fmt.Sprint(child.Id) + " "
	}
	childlist = strings.TrimSpace(childlist)
	ret.Title += " " + childlist

	return
}
func DeleteTodo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	fmt.Println(id)
	todo := models.Todo{Id: id}

	if err := todo.DelTodo(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"Message": err.Error()})

	} else {
		c.JSON(http.StatusOK, todo)
	}

}
func GetTodolist(c *gin.Context) {
	var listquery models.Pageination
	c.DefaultQuery("limit", "10")
	c.DefaultQuery("page", "1")
	limit := c.Query("limit")
	page := c.Query("page")
	if listquery.Limit, _ = strconv.Atoi(limit); listquery.Limit == 0 {
		listquery.Limit = 10
	}
	if listquery.Page, _ = strconv.Atoi(page); listquery.Page == 0 {
		listquery.Page = 1
	}
	listquery.Order = c.Query("order")

	//fmt.Println(listquery)
	todolist, _ := listquery.Listup()
	todliststr := make([]*ViewStruct, 0)
	for _, todo := range todolist.Todolist {
		todliststr = append(todliststr, ModelToView(todo))

	}

	res := TodoListStruct{Todolist: todliststr, Totaldata: todolist.Totaldata, Limit: todolist.Limit, Page: todolist.Page}

	c.JSON(http.StatusOK, res)

}

func (todo *Todo) ConvertToModel() models.Todo {
	ret := models.Todo{Id: todo.Id, Title: todo.Title, Done: todo.Done}
	for _, child := range todo.Children {
		ret.Children = append(ret.Children, &models.Todo{Id: child})
	}
	return ret
}
func CreateTodo(c *gin.Context) {
	var json Todo
	if c.Request.Method == "POST" {
		c.BindJSON(&json)
		todo := json.ConvertToModel()
		if err := todo.CreateTodo(); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusConflict, gin.H{"Message": err.Error()})

		} else {
			todo.FindById()
			c.JSON(http.StatusOK, ModelToView(&todo))
		}

	}

}
func GetTodo(c *gin.Context) {

	//fmt.Println("ttt")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	fmt.Println(id)
	Q := models.Todo{Id: id}
	if err := Q.FindById(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"Message": err.Error()})

	} else {
		c.JSON(http.StatusOK, ModelToView(&Q))
	}

}

//업데이트 요청시 조건을 처리하는 함수
func beforupdatefill(todo *models.Todo) {
	fill := models.Todo{Id: todo.Id}
	err := fill.FindById()
	fmt.Println(err)
	//할일 이름이 비워 있을경우 디비 정보와 일치 시킨다.
	if todo.Title == "" {
		todo.Title = fill.Title
	}
	//참조 하는 부분이 비워져 있을 경우 디비 정보와 일치 시킨다.
	if len(todo.Children) == 0 {
		todo.Children = fill.Children

	}
	//참조하는 값이 0이면 전체 삭제할 수 있도록 한다. 빈 참조리스트를 선언한다.
	if len(todo.Children) == 1 && todo.Children[0].Id == 0 {
		todo.Children = make([]*models.Todo, 0)

	}
	if todo.Done == "" {
		todo.Done = fill.Done
	}

}

//tilte이름 비어 있을 때 처리해야된다.
func ModfiyTodo(c *gin.Context) {
	var json Todo

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if c.Request.Method == "PUT" {
		fmt.Println("run")
		c.BindJSON(&json)
		todo := json.ConvertToModel()
		todo.Id = id
		beforupdatefill(&todo)
		fmt.Println(todo)
		if err := todo.UpdateTodo(); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusConflict, gin.H{"Message": err.Error()})
		} else {
			todo.FindById()
			c.JSON(http.StatusOK, ModelToView(&todo))
		}
	}

}
