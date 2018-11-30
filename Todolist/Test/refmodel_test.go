package Test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"wastest/Todolist/models"
	"wastest/common"
)

func Init() {
	db := common.TestDBInit()
	//db.DropTableIfExists(&models.Todo{})
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

func Makeref(G map[uint][]uint) []*models.Todo {

	ret := make([]*models.Todo, len(G))

	i := 0
	for key, value := range G {
		ret[i] = &models.Todo{
			Id: key,
		}
		for _, ref := range value {

			ret[i].Reflist = append(ret[i].Reflist, &models.Todo{Id: ref})
		}
		i++

	}
	return ret

}
func GetTestCaseList(path string) []string {
	var files []string

	root := path

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return files
}

func InsertTestreflist(G []*models.Todo) {

	db := common.GetDB()
	for i := 0; i < len(G); i++ {
		reflist := G[i].Reflist
		db.Model(&models.Todo{Id: G[i].Id}).Association("Reflist").Append(&reflist)
	}
}

func TestVaildRef(t *testing.T) {
	Init()

	db := common.GetDB()
	i := 0

	for _, filename := range GetTestCaseList("./TestCase/ref/in") {
		if i == 0 {
			i++
			continue
		}

		fi, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		var n int
		fmt.Fscanf(fi, "%d", &n)
		G := make(map[uint][]uint)
		for i := 0; i < n; i++ {
			var u, v uint
			fmt.Fscanf(fi, "%d %d", &u, &v)
			G[u] = append(G[u], v)
		}
		for key, value := range G {
			fmt.Println("Key:", key, "Value:", value)
		}
		var p_id, c_id uint
		fmt.Fscanf(fi, "%d %d", &p_id, &c_id)
		fi.Close()
		refs := Makeref(G)
		InsertTestreflist(refs)
		//val,err :=models.VaildIntersect(p_id,c_id)
		//fmt.Println(val)
		db.Table("ref").Delete("*")
		i++

	}

	//

	/*Case 2
	1->2
	2->3
	4->1
	Try
	3->4
	Answer
	Valid Faild
	*/

}
func TestXx(t *testing.T) {
	//Init()
	u := models.Todo{Id: 1}
	v := models.Todo{Id: 2}
	models.ValidationRef(u, v)

	//InsertDumyData()

}
