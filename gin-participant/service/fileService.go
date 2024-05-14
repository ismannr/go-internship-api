package service

import (
	"gin-crud/initializers"
	model "gin-crud/models"
	"gin-crud/request"
	"gin-crud/response"
	"gin-crud/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func CvUploader(c *gin.Context, participant *model.ParticipantData) {
	var req request.FileRequest

	if err := c.Bind(&req); err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}

	if len(req.FileExt) == 0 {
		response.GlobalResponse(c, "Failed to retrieve file extension", http.StatusBadRequest, nil)
		return
	}

	if !utils.IsPDF(req.FileExt) {
		response.GlobalResponse(c, "Wrong file format", http.StatusBadRequest, nil)
		return
	}

	fileName := utils.GenerateUniqueFileName(req.FileExt)

	cfg, err := getConfig()
	if err != nil {
		response.GlobalResponse(c, "Failed to create cloud session", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	presigner := NewPresigner(cfg)
	key := participant.Cv

	if len(key) == 0 {
		key = "participant/cv/" + fileName
	}

	url, err := presigner.PutObject(os.Getenv("AWS_S3_BUCKET_NAME"), key, 300)
	if err != nil {
		response.GlobalResponse(c, "Failed to update cv", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	if len(participant.Cv) == 0 {
		participant.Cv = key
	}

	if err := initializers.DB.Save(&participant).Error; err != nil {
		response.GlobalResponse(c, "Failed to update participant cv", http.StatusInternalServerError, nil)
		return
	}

	response.GlobalResponse(c, "Successfully generating presigned URL, valid for 5 minute", http.StatusOK, url)
}

func ProfilePictureUploader(c *gin.Context, participant *model.ParticipantData) {
	var req request.FileRequest

	if err := c.Bind(&req); err != nil {
		response.GlobalResponse(c, "Failed to retrieve user request", http.StatusBadRequest, nil)
		return
	}

	if len(req.FileExt) == 0 {
		response.GlobalResponse(c, "Failed to retrieve file extension", http.StatusBadRequest, nil)
		return
	}

	if !utils.IsImage(req.FileExt) {
		response.GlobalResponse(c, "Wrong image format", http.StatusBadRequest, nil)
		return
	}

	fileName := utils.GenerateUniqueFileName(req.FileExt)

	cfg, err := getConfig()
	if err != nil {
		response.GlobalResponse(c, "Failed to create cloud session", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	presigner := NewPresigner(cfg)
	key := participant.ProfilePicture

	if len(key) == 0 {
		key = "participant/profile-picture/" + fileName
	}

	url, err := presigner.PutObject(os.Getenv("AWS_S3_BUCKET_NAME"), key, 300)
	if err != nil {
		response.GlobalResponse(c, "Failed to update profile picture", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	if len(participant.ProfilePicture) == 0 {
		participant.ProfilePicture = key
	}

	if err := initializers.DB.Save(&participant).Error; err != nil {
		response.GlobalResponse(c, "Failed to update participant profile picture", http.StatusInternalServerError, nil)
		return
	}

	response.GlobalResponse(c, "Successfully generating presigned URL, valid for 5 Minute", http.StatusOK, url)
}

func DeleteCv(c *gin.Context, participant *model.ParticipantData) {
	if len(participant.Cv) == 0 {
		response.GlobalResponse(c, "No cv found for deletion", http.StatusNotFound, nil)
		return
	}

	cfg, err := getConfig()
	if err != nil {
		response.GlobalResponse(c, "Failed to create cloud session", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	presigner := NewPresigner(cfg)
	key := participant.Cv

	url, err := presigner.DeleteObject(os.Getenv("AWS_S3_BUCKET_NAME"), key, 300)
	if err != nil {
		response.GlobalResponse(c, "Failed to generate pre-signed delete URL", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	participant.Cv = ""
	if err := initializers.DB.Save(&participant).Error; err != nil {
		response.GlobalResponse(c, "Failed to update participant data after cv deletion", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	response.GlobalResponse(c, "CV deleted successfully", http.StatusOK, url)
}

func DeleteProfilePicture(c *gin.Context, participant *model.ParticipantData) {
	if len(participant.ProfilePicture) == 0 {
		response.GlobalResponse(c, "No profile picture found for deletion", http.StatusNotFound, nil)
		return
	}

	cfg, err := getConfig()
	if err != nil {
		response.GlobalResponse(c, "Failed to create cloud session", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	presigner := NewPresigner(cfg)
	key := participant.ProfilePicture

	url, err := presigner.DeleteObject(os.Getenv("AWS_S3_BUCKET_NAME"), key, 300)
	if err != nil {
		response.GlobalResponse(c, "Failed to generate pre-signed delete URL", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	participant.ProfilePicture = ""
	if err := initializers.DB.Save(&participant).Error; err != nil {
		response.GlobalResponse(c, "Failed to update participant data after profile picture deletion", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	response.GlobalResponse(c, "Profile picture deleted successfully", http.StatusOK, url)
}

func GetProfilePictureURL(c *gin.Context, participant *model.ParticipantData) {
	if len(participant.ProfilePicture) == 0 {
		response.GlobalResponse(c, "No profile picture found", http.StatusNotFound, nil)
		return
	}

	cfg, err := getConfig()
	if err != nil {
		response.GlobalResponse(c, "Failed to create cloud session", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	presigner := NewPresigner(cfg)

	key := participant.ProfilePicture

	req, err := presigner.GetObject(os.Getenv("AWS_S3_BUCKET_NAME"), key, 3600)
	if err != nil {
		response.GlobalResponse(c, "Failed to generate URL for profile picture", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	response.GlobalResponse(c, "Image link generated", http.StatusOK, req)
}

func GetCvURL(c *gin.Context, participant *model.ParticipantData) {
	if len(participant.Cv) == 0 {
		response.GlobalResponse(c, "No CV found", http.StatusNotFound, nil)
		return
	}

	cfg, err := getConfig()
	if err != nil {
		response.GlobalResponse(c, "Failed to create cloud session", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	presigner := NewPresigner(cfg)

	key := participant.Cv

	req, err := presigner.GetObject(os.Getenv("AWS_S3_BUCKET_NAME"), key, 3600)
	if err != nil {
		response.GlobalResponse(c, "Failed to generate URL for cv", http.StatusInternalServerError, nil)
		log.Println(err.Error())
		return
	}

	response.GlobalResponse(c, "CV link generated", http.StatusOK, req)
}
