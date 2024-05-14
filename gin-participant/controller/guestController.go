package controller

import (
	"gin-crud/config"
	"gin-crud/service"
	"github.com/gin-gonic/gin"
)

func GuestController(r *gin.Engine) {
	r.POST("/sign-up", service.UserRegister)
	r.POST("/login", service.Login)
	//r.GET("/provinces", service.GetProvinceList)
	//r.GET("/cities/:id", service.GetCityDependsOnProvince)
	//r.GET("/city/:id", service.GetCity)
	r.GET("/logout", config.AuthFilter, service.Logout)
	r.POST("/forgot-password", service.RecoveryPassword)
	r.PUT("/reset-password", service.RecoveryPassword)
	r.POST("/reset-password/verifying-token/:token", config.RecoveryAuthFilter, service.VerifyResetPasswordToken)
	r.POST("/reset-password/change-password/:token", config.RecoveryAuthFilter, service.ResetPassword)

	r.GET("/email-confirmation/:token", config.EmailConfirmFilter, service.ConfirmingEmail)
	r.GET("/resend-email-confirmation/:email", config.ResendConfirmFilter, service.ResendEmailConfirm)
}
