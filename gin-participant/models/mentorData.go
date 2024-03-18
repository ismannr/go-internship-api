package models

import "gorm.io/gorm"

type MentorData struct {
	gorm.Model
	MentorId          int64      `json:"mentor_id"`
	MentorName        string     `json:"mentor_name"`
	MentorPicture     []byte     `json:"mentor_picture"`
	MentorEmployeeId  string     `json:"mentor_employee_id"`
	Email             string     `json:"email"`
	MentorPhoneNumber string     `json:"mentor_phone_number"`
	MentorPosition    string     `json:"mentor_position"`
	SystemDataID      uint       `gorm:"column:system_data_id"`
	SystemData        SystemData `gorm:"foreignKey:SystemDataID"`
}
