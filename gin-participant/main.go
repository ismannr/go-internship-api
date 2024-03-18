package main

import (
	"gin-crud/controller"
	"gin-crud/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DatabaseInit()
}

func main() {
	r := gin.Default()
	controller.UserController(r)
	controller.GuestController(r)
	controller.AdminController(r)
	r.Run() // listen and serve on default 0.0.0.0:8080
}
