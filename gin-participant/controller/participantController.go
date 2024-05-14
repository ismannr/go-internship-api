package controller

import (
	"gin-crud/config"
	"gin-crud/service"
	"github.com/gin-gonic/gin"
)

func UserController(r *gin.Engine) {
	r.GET("/user", config.AuthFilter, service.GetParticipantData)
	r.PUT("/user", config.AuthFilter, service.UpdateParticipant)

	r.PUT("/user/cv", config.AuthFilter, service.UploadCvByAuth)
	r.PUT("/user/profile-picture", config.AuthFilter, service.UploadProfilePictureByAuth)

	r.DELETE("/user/cv/delete", config.AuthFilter, service.DeleteCvByAuth)
	r.DELETE("/user/profile-picture/delete", config.AuthFilter, service.DeleteProfilePictureByAuth)

	r.GET("/user/profile-picture", config.AuthFilter, service.GetProfilePictureByAuth)
	r.GET("/user/cv", config.AuthFilter, service.GetCvByAuth)
}
