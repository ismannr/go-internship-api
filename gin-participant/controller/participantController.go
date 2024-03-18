package controller

import (
	"gin-crud/config"
	"gin-crud/service"
	"github.com/gin-gonic/gin"
)

func UserController(r *gin.Engine) {
	r.GET("/user", config.AuthFilter, service.GetParticipantData)
	r.PUT("/user", config.AuthFilter, service.UpdateParticipant)

}
