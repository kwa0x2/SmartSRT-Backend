package domain

import (
	"mime/multipart"
)

type LambdaBodyResponse struct {
	Message string `json:"message"`
	SRTURL  string `json:"srt_url"`
}

type LambdaResponse struct {
	StatusCode int                `json:"status_code"`
	Body       LambdaBodyResponse `json:"body"`
}

type FileConversionRequest struct {
	UserID       string               `json:"user_id"`
	WordsPerLine int                  `json:"words_per_line"`
	FileName     string               `json:"file_name"`
	File         multipart.File       `json:"file"`
	FileHeader   multipart.FileHeader `json:"file_header"`
}

type SubtitleUseCase interface {
	UploadFileAndConvertToSRT(request FileConversionRequest) (*LambdaResponse, error)
}

type SubtitleRepository interface {
	UploadFileToS3(request FileConversionRequest) (string, error)
	TriggerLambdaFunc(request FileConversionRequest) (*LambdaResponse, error)
}
