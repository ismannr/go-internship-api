package config

import (
	"fmt"
	"gin-crud/initializers"
	model "gin-crud/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func AuthFilter(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user model.ParticipantData
		if err := initializers.DB.First(&user, claims["sub"]).Error; err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)

	}
}

func AdminAuthFilter(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var user model.ParticipantData
		if err := initializers.DB.Preload("SystemData").First(&user, claims["sub"]).Error; err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if user.SystemData.Level != model.LevelAdmin {
			fmt.Println("test")
			c.AbortWithStatus(403)
			return
		}
		c.Set("user", user)
		c.Next()
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
