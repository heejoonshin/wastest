package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"wastest/Todolist/Route"
	"wastest/Todolist/models"
	"wastest/common"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Todo{})
}

var router *gin.Engine

func main() {
	db := common.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()

	v1 := r.Group("/api")
	Route.TodoGroup(v1.Group("/todo"))
	Route.TodolistGroup(v1.Group("/todolist"))
	r.Run()
}
