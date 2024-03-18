package controller

import (
	"gin-crud/config"
	"gin-crud/service"
	"github.com/gin-gonic/gin"
)

func AdminController(r *gin.Engine) {
	r.GET("/admin/user/list", config.AdminAuthFilter, service.GetParticipantList)
	r.GET("/admin/user/:id", config.AdminAuthFilter, service.GetParticipantById)
	r.DELETE("/admin/user/:id", config.AdminAuthFilter, service.DeleteUserById)
	r.PUT("/admin/user/:id", config.AdminAuthFilter, service.UpdateParticipantById)
	r.POST("/admin/user", config.AdminAuthFilter, service.GetParticipantByEmail)
	r.POST("/admin/create-user", config.AdminAuthFilter, service.CreateParticipant)
}
