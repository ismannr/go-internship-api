package controller

import (
	"gin-crud/config"
	"gin-crud/service"
	"github.com/gin-gonic/gin"
)

func GuestController(r *gin.Engine) {
	r.POST("/sign-up", service.UserRegister)
	r.POST("/login", service.Login)
	r.GET("/provinces", service.GetProvinceList)
	r.GET("/cities/:id", service.GetCityDependsOnProvince)
	r.GET("/city/:id", service.GetCity)
	r.GET("/logout", config.AuthFilter, service.Logout)
}
