package delivery

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain/types"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/middleware"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/rabbitmq"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"github.com/kwa0x2/AutoSRT-Backend/utils/validator"
)

type SRTDelivery struct {
	SRTUseCase domain.SRTUseCase
	RabbitMQ   *domain.RabbitMQ
}

func (sd *SRTDelivery) ConvertFileToSRT(ctx *gin.Context) {
	startTime := time.Now()

	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userData := user.(*domain.User)

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "file_form_parse"})
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("File is required. Please try again."))
		return
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "file_close"})
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
		utils.HandleErrorWithSentry(ctx, fmt.Errorf("file does not implement io.Seeker interface"), map[string]interface{}{"action": "file_seeker_check"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to process file. Please try again."))
		return
	}

	duration, err := utils.GetMediaDuration(file, fileType)
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "get_media_duration"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to get file duration. Please try again."))
		return
	}

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
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "file_seek_reset"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to process file. Please try again."))
		return
	}

	params, err := validator.ValidateConversionParams(ctx)
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "validate_conversion_params"})
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse(err.Error()))
		return
	}

	fileID := utils.GenerateUUID()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "file_read_bytes"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to read file. Please try again."))
		return
	}

	msg := domain.ConversionMessage{
		UserID:              userData.ID,
		WordsPerLine:        params.WordsPerLine,
		Punctuation:         params.Punctuation,
		ConsiderPunctuation: params.ConsiderPunctuation,
		FileID:              fileID,
		FileName:            header.Filename,
		FileContent:         fileBytes,
		FileSize:            header.Size,
		FileDuration:        duration,
		Email:               userData.Email,
	}

	response, err := rabbitmq.PublishConversionMessage(sd.RabbitMQ, ctx, msg)
	if err != nil {
		if err.Error() == "Response timeout" {
			ctx.JSON(http.StatusAccepted, gin.H{
				"message": "Your file is being processed. You will receive an email when it's ready.",
				"file_id": fileID,
			})
			return
		}
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "rabbitmq_publish"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to queue conversion. Please try again."))
		return
	}

	middleware.RecordSRTMetrics("queued_success", time.Since(startTime))
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
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "find_user_histories"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred while retrieving history data. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, srtHistoriesData)
}

func (sd *SRTDelivery) TestSentry(ctx *gin.Context) {
	utils.HandleErrorWithSentry(ctx, fmt.Errorf("test error2ssss"), map[string]interface{}{"state": "test22"})

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Test error sent to Sentry successfully",
		"timestamp": time.Now().Unix(),
		"user_id":   "anonymous",
	})
}
