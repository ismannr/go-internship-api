package service

import (
	"gin-crud/initializers"
	model "gin-crud/models"
	"gin-crud/request"
	"gin-crud/response"
	"gin-crud/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

func Login(c *gin.Context) {
	var user model.SystemData
	var req struct {
		Email    string
		Password string
	}
	if err := c.Bind(&req); err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}
	initializers.DB.First(&user, "email = ?", req.Email)
	if user.ID == 0 {
		response.GlobalResponse(c, "Invalid email or password", http.StatusBadRequest, nil)
		return
	}
	if utils.HashIsMatched(user.Password, req.Password) == false {
		response.GlobalResponse(c, "Invalid email or password", http.StatusBadRequest, nil)
		return
	}
	generateToken(user, c)
}

func UserRegister(c *gin.Context) {
	var req request.ParticipantRequest
	var s strings.Builder
	var userDB model.ParticipantData
	var isSatisfied bool = true

	if err := c.Bind(&req); err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}

	if len(req.MentorId) == 0 || (!regexp.MustCompile(`^M\d{3}$`).MatchString(req.MentorId) && len(req.MentorId) != 4) {
		s.WriteString("Mentor ID, ")
		isSatisfied = false
	}
	if len(req.Name) < 3 {
		s.WriteString("Name, ")
		isSatisfied = false
	}
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		s.WriteString("Email (wrong format), ")
		isSatisfied = false
	}
	if exist := initializers.DB.Where("email = ?", req.Email).First(&userDB).Error; exist == nil {
		s.WriteString("Email already exist, ")
		isSatisfied = false
	}
	if strings.ToLower(req.Gender) != "male" && strings.ToLower(req.Gender) != "female" {
		s.WriteString("Gender, ")
		isSatisfied = false
	}
	dob, err := utils.ParseDate(req.Dob)
	if err != nil {
		s.WriteString("DOB (wrong date format), ")
		isSatisfied = false
	} else {
		if !utils.IsAdult(req.Dob) {
			s.WriteString("DOB (Must be over 17), ")
			isSatisfied = false
		}
	}

	if len(req.Password) < 8 {
		s.WriteString("Password min. 8 char, ")
		isSatisfied = false
	} else {
		uppercaseRegex := regexp.MustCompile(`[A-Z]`)
		numberRegex := regexp.MustCompile(`[0-9]`)
		if !uppercaseRegex.MatchString(req.Password) || !numberRegex.MatchString(req.Password) {
			s.WriteString("Password (Must contain at least one uppercase letter and one number), ")
			isSatisfied = false
		} else if req.ConfirmPass != req.Password {
			s.WriteString("Confirmation Password doesn't match, ")
			isSatisfied = false
		}
	}
	if exist := initializers.DB.Where("phone_number = ?", req.PhoneNumber).First(&userDB).Error; exist == nil {
		s.WriteString("Phone number already exist, ")
		isSatisfied = false
	}

	if !isSatisfied {
		message := "User data requirements not satisfied: " + s.String()
		response.GlobalResponse(c, message, http.StatusBadRequest, nil)
		return
	}

	password, err := utils.HashEncoder(req.Password)
	if err != nil {
		message := "Error encoding the password"
		response.GlobalResponse(c, message, http.StatusBadRequest, nil)
		return
	}

	systemUser := model.SystemData{
		Email:    req.Email,
		Password: password,
		Role:     model.RoleParticipant,
		Level:    model.LevelUser,
	}
	user := model.ParticipantData{
		Name:         req.Name,
		MentorId:     strings.ToUpper(req.MentorId),
		Gender:       strings.ToUpper(req.Gender),
		Dob:          dob,
		Email:        req.Email,
		University:   req.University,
		Gpa:          req.Gpa,
		Address:      req.Address,
		Province:     req.Province,
		City:         req.City,
		PhoneNumber:  req.PhoneNumber,
		SystemDataID: systemUser.ID,
		SystemData:   systemUser,
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		response.GlobalResponse(c, "Failed to save user data", http.StatusBadRequest, nil)
		return
	}
	initializers.DB.Save(&user)
	response.GlobalResponse(c, "Participant_data created successfully", http.StatusOK, user)
}

func MentorRegister(c *gin.Context) {

}
