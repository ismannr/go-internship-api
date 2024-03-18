package service

import (
	"errors"
	"gin-crud/initializers"
	models "gin-crud/models"
	"gin-crud/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func generateToken(user models.SystemData, c *gin.Context) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	switch user.Role {
	case models.RoleMentor:
		claims["sub"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour).Unix()
		claims["role"] = string(user.Role)
	case models.RoleParticipant:
		claims["sub"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour).Unix()
		claims["role"] = string(user.Role)
	default:
		response.GlobalResponse(c, "Invalid user type", http.StatusInternalServerError, nil)
		return
	}

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		response.GlobalResponse(c, "Invalid token creation", http.StatusInternalServerError, nil)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, int(time.Hour.Seconds()), "", "", false, true)
	response.GlobalResponse(c, "Token generated", 200, nil)
}

func getUserByAuth(c *gin.Context) (interface{}, error) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		response.GlobalResponse(c, "Authorization token not provided", http.StatusUnauthorized, nil)
		return nil, errors.New("authorization token not provided")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["sub"].(float64))

		if claims["role"] == string(models.RoleMentor) {
			var mentor models.MentorData
			if err := initializers.DB.Preload("SystemData").First(&mentor, "id = ?", userID).Error; err == nil {
				return &mentor, nil
			}
		} else if claims["role"] == string(models.RoleParticipant) {
			var participant models.ParticipantData
			if err := initializers.DB.Preload("SystemData").First(&participant, "id = ?", userID).Error; err == nil {
				return &participant, nil
			}
		}
		return nil, errors.New("user not found")
	}
	return nil, errors.New("invalid token")
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	response.GlobalResponse(c, "Logout successful", 200, nil)
}
