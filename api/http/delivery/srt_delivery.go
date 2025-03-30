package delivery

import (
	"errors"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"io"
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

	sessionUserRole, exists := ctx.Get("role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Unauthorized. Please log in and try again."))
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("File is required. Please try again."))
		return
	}
	defer file.Close()

	fileType := filepath.Ext(header.Filename)

	if sessionUserRole != types.Pro && fileType == ".wav" {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("You need to upgrade to the Pro plan to upload WAV files."))
		return
	}

	if !utils.IsValidMediaFile(fileType) {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid file format. Only mp4, mp3 and wav files are accepted."))
		return
	}

	seeker, ok := file.(io.Seeker)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to process file. Please try again."))
		return
	}

	duration, err := utils.GetMediaDuration(file, fileType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to get file duration. Please try again."))
		return
	}

	_, err = seeker.Seek(0, io.SeekStart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to process file. Please try again."))
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
		FileDuration:        duration,
	}

	response, err := sd.SRTUseCase.UploadFileAndConvertToSRT(request)
	if err != nil {
		if errors.Is(err, utils.ErrLimitReached) {
			ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse(
				"Monthly limit reached. Upgrade to Premium for more conversions.",
			))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(
			"Failed to generate SRT file. Please check your file and try again.",
		))
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
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred while retrieving history data. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, srtHistoriesData)
}
