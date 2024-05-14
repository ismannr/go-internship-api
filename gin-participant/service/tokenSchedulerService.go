package service

import (
	"errors"
	"gin-crud/initializers"
	model "gin-crud/models"
	"gorm.io/gorm"
	"log"
	"time"
)

func TokenExpirationCheckAndUpdateScheduler() {
	tokenExpirationCheckAndUpdate()
	ticker := time.NewTicker(time.Minute * 30)
	defer ticker.Stop()
	for range ticker.C {
		tokenExpirationCheckAndUpdate()
	}
}

func tokenExpirationCheckAndUpdate() {
	var sysData []model.SystemData

	initializers.DB.Where("last_login < ?", time.Now().Add(-1*time.Hour)).Find(&sysData)

	for _, session := range sysData {
		session.CurrentlyLogin = false

		if err := initializers.DB.Save(&session).Error; err != nil {
			log.Println("Failed to update session:", err)
			continue
		}

		if session.TokenID == nil {
			continue
		}

		var token model.Token
		if err := initializers.DB.First(&token, session.TokenID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			log.Println("Failed to retrieve token data:", err)
			continue
		}

		if err := initializers.DB.Unscoped().Delete(&token).Error; err != nil {
			log.Println("Failed to invalidate token:", err)
			continue
		}
	}
}
