package delivery

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type SRTDelivery struct {
	SRTUseCase domain.SRTUseCase
}

func (sd *SRTDelivery) ConvertFileToSRT(ctx *gin.Context) {
	//sessionUserID, exists := ctx.Get("user_id")
	//if !exists {
	//	ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Unauthorized. Please log in and try again."))
	//	return
	//}
	//
	//userIDStr, ok := sessionUserID.(string)
	//if !ok {
	//	ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid session data. Please log in again. If the issue persists, contact support."))
	//	return
	//}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("File is required. Please try again."))
		return
	}
	defer file.Close()

	if !isValidFile(header.Filename) {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid file format. Only mp4 files are accepted. Please try again."))
		return
	}

	wordsPerLineStr := ctx.PostForm("words_per_line")
	if wordsPerLineStr == "" {
		wordsPerLineStr = "3"
	}

	wordsPerLine, err := strconv.Atoi(wordsPerLineStr)
	if err != nil {
		wordsPerLine = 3
	}

	if wordsPerLine < 1 || wordsPerLine > 5 {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Words per line must be between 1 and 5."))
		return
	}

	request := domain.FileConversionRequest{
		UserID:       "userIDStr",
		WordsPerLine: wordsPerLine,
		File:         file,
		FileHeader:   *header,
	}

	response, err := sd.SRTUseCase.UploadFileAndConvertToSRT(request)
	if err != nil {
		//ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response)
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
