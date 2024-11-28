package api

import (
	"errors"
	"net/http"
	db "ocr/db/sqlc"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tiagomelo/go-ocr/ocr"
)

type CreateImageConversionRequest struct {
	UserID        int32  `form:"user_id"`
	ImageName     string `json:"image_name"`
	ExtractedText string `json:"extracted_text"`
}

type CreateImageConversionResponse struct {
	UserID        int32  `form:"user_id"`
	ImageName     string `json:"image_name"`
	ExtractedText string `json:"extracted_text"`
}

func (server Server) CreateImageConversion(ctx *gin.Context) {
	var req CreateImageConversionRequest

	if err := ctx.ShouldBind(&req); err != nil {
		err = errors.New("input is not valid, Please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file"})
		return
	}

	// Save the file temporarily
	savePath := filepath.Join("uploads", file.Filename)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Extract text using Tesseract OCR
	t, err := ocr.New()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Tesseract"})
		return
	}

	// Extract text from the uploaded image
	extractedText, err := t.TextFromImageFile(savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract text"})
		return
	}
	defer os.Remove(savePath)

	arg := db.CreateImageConversionParams{
		UserID:        req.UserID,
		ImageName:     file.Filename,
		ExtractedText: strings.TrimSpace(extractedText),
	}

	image, err := server.store.CreateImageConversion(ctx, arg)
	if err != nil {
		// err = errors.New("failed to upload profile picture, Please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, image)
}

type GetImageConversionRequest struct {
	UserID        int32  `json:"user_id"`
	ImageName     string `json:"image_name"`
	ExtractedText string `json:"extracted_text"`
}

type GetImageConversionResponse struct {
	UserID        int32  `json:"user_id"`
	ImageName     string `json:"image_name"`
	ExtractedText string `json:"extracted_text"`
}

func (server Server) GetImageConversion(ctx *gin.Context) {

}
