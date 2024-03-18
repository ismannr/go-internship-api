package request

import "gin-crud/models"

type ParticipantRequest struct {
	MentorId       string       `json:"mentor_id"`
	Name           string       `json:"participant_name"`
	Email          string       `json:"email" gorm:"unique"`
	Password       string       `json:"password"`
	ConfirmPass    string       `json:"confirm_password"`
	Gender         string       `json:"gender"`
	PhoneNumber    string       `json:"participant_phone_number" gorm:"unique"`
	Dob            string       `json:"participant_birth_date"`
	University     string       `json:"participant_university"`
	Address        string       `json:"participant_address"`
	Gpa            float64      `json:"participant_gpa"`
	Cv             string       `json:"participant_cv"`
	ProfilePicture string       `json:"profile_picture"`
	Status         int64        `json:"participant_status"`
	City           string       `json:"participant_city"`
	Province       string       `json:"participant_province"`
	Role           models.Role  `json:"role"`
	Level          models.Level `json:"level"`
}
