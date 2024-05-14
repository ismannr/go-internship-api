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
	r.PUT("/admin/user/activation-status/:id", config.AdminAuthFilter, service.IsActivated)

	r.GET("/admin/user/cv/:id", config.AdminAuthFilter, service.GetCvById)
	r.DELETE("/admin/user/cv/:id", config.AdminAuthFilter, service.DeleteCvById)
	r.PUT("/admin/user/cv/:id", config.AdminAuthFilter, service.UploadCvById)

	r.GET("/admin/user/profile-picture/:id", config.AdminAuthFilter, service.GetProfilePictureById)
	r.DELETE("/admin/user/profile-picture/:id", config.AdminAuthFilter, service.DeleteProfilePictureById)
	r.PUT("/admin/user/profile-picture/:id", config.AdminAuthFilter, service.UploadProfilePictureById)

	r.GET("/admin/user/email/:email", config.AdminAuthFilter, service.GetParticipantByEmail)
	r.POST("/admin/create-user", config.AdminAuthFilter, service.CreateParticipant)

}
