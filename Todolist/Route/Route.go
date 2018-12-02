package Route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"wastest/Todolist/models"
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

func ModelToView(i interface{}) {
	rt := reflect.TypeOf(i)
	if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
		fmt.Println(rt.Elem())
	} else {
		Values := reflect.ValueOf(i).Elem()

		fmt.Println(Values)
	}

}

func TodoGroup(router *gin.RouterGroup) {
	router.POST("/", CreateTodo)
	router.GET("/:id", GetTodo)
	router.PUT("/:id", ModfiyTodo)
	router.DELETE("/:id", DeleteTodo)
	//router.GET("/list", GetTodolist)

}

func TodolistGroup(router *gin.RouterGroup) {

	router.GET("/", GetTodolist)

}
func DeleteTodo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	fmt.Println(id)
	todo := models.Todo{Id: id}

	todo.DelTodo()
	c.JSON(http.StatusOK, todo)

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

	c.JSON(http.StatusOK, todolist)

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

		}
		c.JSON(http.StatusOK, todo)
	}

}
func GetTodo(c *gin.Context) {

	//fmt.Println("ttt")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	fmt.Println(id)
	Q := models.Todo{Id: id}
	Q.FindById()
	c.JSON(http.StatusOK, Q)

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
		fmt.Println(todo)
		if err := todo.UpdateTodo(); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusConflict, gin.H{"Message": err.Error()})
		} else {
			c.JSON(http.StatusOK, todo)
		}
	}

}
