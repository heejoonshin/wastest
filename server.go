package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"wastest/Todolist/models"
	"wastest/common"

	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Todo{})
}

func main() {

	db := common.Init()
	Migrate(db)
	defer db.Close()

	//var t []*Todolist.Todo

	//db.Model(&t).Association("Reflist").Append(&Todolist.Todo{Id: 3}, &Todolist.Todo{Id: 4})
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./views", true)))
	page := models.Pageination{
		Limit: 10,
		Order: "done",
		Page:  1,
	}
	f, _ := page.Listup()
	fmt.Println(f[0])

	/*
		v1 := r.Group("/api")
		users.UsersRegister(v1.Group("/users"))
		v1.Use(users.AuthMiddleware(false))
		articles.ArticlesAnonymousRegister(v1.Group("/articles"))
		articles.TagsAnonymousRegister(v1.Group("/tags"))

		v1.Use(users.AuthMiddleware(true))
		users.UserRegister(v1.Group("/user"))
		users.ProfileRegister(v1.Group("/profiles"))

		articles.ArticlesRegister(v1.Group("/articles"))

		testAuth := r.Group("/api/ping")

		testAuth.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// test 1 to 1
		tx1 := db.Begin()
		userA := users.UserModel{
			Username: "AAAAAAAAAAAAAAAA",
			Email:    "aaaa@g.cn",
			Bio:      "hehddeda",
			Image:    nil,
		}
		tx1.Save(&userA)
		tx1.Commit()
		fmt.Println(userA)

		//db.Save(&ArticleUserModel{
		//    UserModelID:userA.ID,
		//})
		//var userAA ArticleUserModel
		//db.Where(&ArticleUserModel{
		//    UserModelID:userA.ID,
		//}).First(&userAA)
		//fmt.Println(userAA)
	*/
	r.Run() // listen and serve on 0.0.0.0:8080
}
