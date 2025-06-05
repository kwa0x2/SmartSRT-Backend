package delivery

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain/types"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"github.com/kwa0x2/AutoSRT-Backend/utils/validator"
)

type SRTDelivery struct {
	SRTUseCase domain.SRTUseCase
}

func (sd *SRTDelivery) ConvertFileToSRT(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userData := user.(*domain.User)

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("File is required. Please try again."))
		return
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		}
	}(file)

	fileType := filepath.Ext(header.Filename)

	if userData.Plan != types.Pro && fileType == ".wav" {
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

	// Dosya süresi kontrolü
	maxDuration := 30 * time.Second
	if userData.Plan == types.Pro {
		maxDuration = 5 * time.Minute
	}

	fileDuration := time.Duration(duration * float64(time.Second))
	if fileDuration > maxDuration {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse(
			"File duration exceeds the limit. Maximum duration is "+maxDuration.String()+" for your plan.",
		))
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

	request := domain.FileConversionRequest{
		UserID:              userData.ID,
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
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (sd *SRTDelivery) FindHistories(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userData := user.(*domain.User)

	srtHistoriesData, err := sd.SRTUseCase.FindHistoriesByUserID(userData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred while retrieving history data. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, srtHistoriesData)
}
