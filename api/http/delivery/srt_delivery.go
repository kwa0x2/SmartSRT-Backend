package delivery

import (
	"fmt"
	"net/http"
	"path/filepath"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"github.com/kwa0x2/AutoSRT-Backend/utils/validator"
)

type SRTDelivery struct {
	SRTUseCase domain.SRTUseCase
}

func (sd *SRTDelivery) ConvertFileToSRT(ctx *gin.Context) {
	sessionUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Unauthorized. Please log in and try again."))
		return
	}

	userIDStr, ok := sessionUserID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid session data. Please log in again. If the issue persists, contact support."))
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("File is required. Please try again."))
		return
	}
	defer file.Close()

	if !isValidFile(header.Filename) {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid file format. Only mp4 files are accepted."))
		return
	}

	params, err := validator.ValidateConversionParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse(err.Error()))
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	request := domain.FileConversionRequest{
		UserID:              userID,
		WordsPerLine:        params.WordsPerLine,
		Punctuation:         params.Punctuation,
		ConsiderPunctuation: params.ConsiderPunctuation,
		File:                file,
		FileHeader:          *header,
	}

	response, err := sd.SRTUseCase.UploadFileAndConvertToSRT(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(fmt.Sprintf("Conversion failed: %v", err)))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (sd *SRTDelivery) FindHistories(ctx *gin.Context) {
	sessionUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Unauthorized. Please log in and try again."))
		return
	}

	userIDStr, ok := sessionUserID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid session data. Please log in again. If the issue persists, contact support."))
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	srtHistoriesData, err := sd.SRTUseCase.FindHistoriesByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, srtHistoriesData)
}

func isValidFile(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp4":
		return true
	default:
		return false
	}
}
