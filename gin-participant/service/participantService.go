package service

import (
	"errors"
	"fmt"
	"gin-crud/utils"
	"gorm.io/gorm"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"gin-crud/initializers"
	model "gin-crud/models"
	"gin-crud/request"
	"gin-crud/response"
	"github.com/gin-gonic/gin"
)

func getParticipantByAuth(c *gin.Context) (*model.ParticipantData, error) {
	user, err := getUserByAuth(c)
	if err != nil {
		return nil, err
	}

	participant, ok := user.(*model.ParticipantData)
	if !ok {
		return nil, errors.New("invalid participant data")
	}
	return participant, nil
}

func GetParticipantData(c *gin.Context) {
	participant, err := getParticipantByAuth(c)
	if err != nil {
		response.GlobalResponse(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}
	response.GlobalResponse(c, "Successfully retrieving participant data", http.StatusOK, participant)
}

func UpdateParticipant(c *gin.Context) {
	participant, err := getParticipantByAuth(c)
	if err != nil {
		response.GlobalResponse(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}
	var req request.ParticipantRequest
	if err := c.Bind(&req); err != nil {
		response.GlobalResponse(c, "Error binding the requested data", http.StatusBadRequest, err)
		return
	}

	message, err, status, participant := validateParticipantRequest(req, participant)
	if err != nil || status != 200 {
		response.GlobalResponse(c, message, status, nil)
		return
	}

	if err := initializers.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&participant).Error; err != nil {
		response.GlobalResponse(c, "Failed to update participant data", http.StatusInternalServerError, nil)
		return
	}
	response.GlobalResponse(c, message, http.StatusOK, participant)

}

func validateParticipantRequest(req request.ParticipantRequest, participant *model.ParticipantData) (string, error, int, *model.ParticipantData) {
	var invalid []string
	var valid []string
	var isSatisfied bool = true

	if len(req.Name) != 0 && req.Name != participant.MentorId {
		if len(req.Name) < 3 {
			invalid = append(invalid, fmt.Sprintf("Name(min. 3 char)"))
			isSatisfied = false
		} else if participant.Name != req.Name && len(req.Name) >= 3 {
			valid = append(valid, fmt.Sprintf("Name"))
			participant.Name = req.Name
		}
	}

	if len(req.Password) != 0 {
		if len(req.Password) < 8 {
			invalid = append(invalid, "Password must be at least 8 characters")
			isSatisfied = false
		} else if utils.HashIsMatched(participant.SystemData.Password, req.Password) == true {
			invalid = append(invalid, "New password cannot be the same as the previous password")
			isSatisfied = false
		} else {
			uppercaseRegex := regexp.MustCompile(`[A-Z]`)
			numberRegex := regexp.MustCompile(`[0-9]`)
			if !uppercaseRegex.MatchString(req.Password) || !numberRegex.MatchString(req.Password) {
				invalid = append(invalid, "Password must contain at least one uppercase letter and one number")
				isSatisfied = false
			} else {
				hashedPassword, err := utils.HashEncoder(req.Password)
				if err != nil {
					return "Error encoding the password", err, 500, nil
				}
				valid = append(valid, "Password")
				participant.SystemData.Password = hashedPassword
			}
		}
	}

	if len(req.MentorId) != 0 && req.MentorId != participant.MentorId {
		if !regexp.MustCompile(`^M\d{3}$`).MatchString(req.MentorId) && len(req.MentorId) != 4 {
			invalid = append(invalid, "MentorID (must start with 'M' followed by 3 digits of number)")
			isSatisfied = false
		} else {
			valid = append(valid, fmt.Sprintf("Mentor ID"))
			participant.MentorId = strings.ToUpper(req.MentorId)
		}
	}

	if len(req.Email) != 0 && req.Email != participant.Email {
		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			invalid = append(invalid, "Email(wrong format)")
			isSatisfied = false
		} else {
			participant.Email = req.Email
			participant.SystemData.Email = req.Email
			valid = append(valid, fmt.Sprintf("Email"))
		}
	}

	if len(req.Gender) != 0 && strings.ToLower(req.Gender) != strings.ToLower(participant.Gender) {
		if strings.ToLower(req.Gender) != "male" && strings.ToLower(req.Gender) != "female" {
			invalid = append(invalid, "Gender(must be 'male' or 'female')")
			isSatisfied = false
		} else {
			valid = append(valid, "Gender")
			participant.Gender = strings.ToUpper(req.Gender)
		}
	}

	if len(req.PhoneNumber) != 0 && req.PhoneNumber != participant.PhoneNumber {
		if !regexp.MustCompile(`^\d{10,14}$`).MatchString(req.PhoneNumber) {
			invalid = append(invalid, "Phone number (must consist of 10-14 digits)")
			isSatisfied = false
		} else {
			valid = append(valid, "Phone Number")
			participant.PhoneNumber = req.PhoneNumber
		}
	}

	if len(req.University) != 0 && req.University != participant.University {
		if len(req.University) < 5 {
			invalid = append(invalid, "University(min. 5 character)")
			isSatisfied = false
		} else {
			valid = append(valid, "University")
			participant.University = req.University
		}
	}

	if len(req.Address) != 0 && req.Address != participant.Address {
		if len(req.Address) < 5 {
			invalid = append(invalid, "Address(min. 5 character)")
			isSatisfied = false
		} else {
			valid = append(valid, "Address")
			participant.Address = req.Address
		}
	}
	if len(req.Province) != 0 && req.Province != participant.Province {
		valid = append(valid, "Address")
		participant.Address = req.Address

	}
	if len(req.City) != 0 && req.City != participant.City {
		valid = append(valid, "Address")
		participant.Address = req.Address
	}

	if req.Gpa != 0 && req.Gpa != participant.Gpa {
		if req.Gpa < 0.01 || req.Gpa > 4.00 {
			invalid = append(invalid, "GPA(must be between 0.01 and 4.00)")
			isSatisfied = false
		} else {
			valid = append(valid, "GPA")
			participant.Gpa = req.Gpa
		}
	}

	if !isSatisfied {
		return "Invalid fields: " + strings.Join(invalid, ", "), nil, 400, nil
	}
	return "Updated fields: " + strings.Join(valid, ", "), nil, 200, participant
}

func UploadProfilePictureByAuth(c *gin.Context) {
	participant, err := getParticipantByAuth(c)
	if err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, err.Error())
		return
	}
	ProfilePictureUploader(c, participant)
}

func UploadCvByAuth(c *gin.Context) {
	participant, err := getParticipantByAuth(c)
	if err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}
	CvUploader(c, participant)
}

func DeleteProfilePictureByAuth(c *gin.Context) {
	participant, err := getParticipantByAuth(c)
	if err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}
	DeleteProfilePicture(c, participant)
}

func DeleteCvByAuth(c *gin.Context) {
	participant, err := getParticipantByAuth(c)
	if err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}
	DeleteCv(c, participant)
}

func GetProfilePictureByAuth(c *gin.Context) {
	participant, err := getParticipantByAuth(c)
	if err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}
	GetProfilePictureURL(c, participant)
}

func GetCvByAuth(c *gin.Context) {
	participant, err := getParticipantByAuth(c)
	if err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}
	GetCvURL(c, participant)
}
