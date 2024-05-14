package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ParticipantData struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key"`
	gorm.Model
	MentorId       string      `json:"mentor_id"`
	Name           string      `json:"participant_name"`
	Email          string      `json:"email" gorm:"uniqueIndex"`
	Gender         string      `json:"gender"`
	PhoneNumber    string      `json:"participant_phone_number" gorm:"uniqueIndex"`
	Dob            time.Time   `json:"participant_birth_date"`
	University     string      `json:"participant_university"`
	Address        string      `json:"participant_address"`
	Gpa            float64     `json:"participant_gpa"`
	Cv             string      `json:"participant_cv"`
	ProfilePicture string      `json:"profile_picture"`
	Status         int64       `json:"participant_status"`
	City           string      `json:"participant_city"`
	Province       string      `json:"participant_province"`
	SystemDataID   *uuid.UUID  `gorm:"column:system_data_id;uniqueIndex"`
	SystemData     *SystemData `gorm:"foreignKey:SystemDataID;constraint:OnDelete:CASCADE;"`
}
