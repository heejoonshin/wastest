package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"wastest/Todolist/Route"
	"wastest/Todolist/models"
	"wastest/common"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Todo{})
}

func main() {
	db := common.Init()
	Migrate(db)
	defer db.Close()

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./view", true)))

	router.Use(cors.Default())

	v1 := router.Group("/api")
	Route.TodoGroup(v1.Group("/todo"))
	Route.TodolistGroup(v1.Group("/todolist"))

	router.Run()
}
