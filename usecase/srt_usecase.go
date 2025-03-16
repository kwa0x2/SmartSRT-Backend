package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"strings"
	"time"
)

type srtUseCase struct {
	srtRepository domain.SRTRepository
}

func NewSRTUseCase(srtRepository domain.SRTRepository) domain.SRTUseCase {
	return &srtUseCase{
		srtRepository: srtRepository,
	}
}

func (su *srtUseCase) UploadFileAndConvertToSRT(request domain.FileConversionRequest) (*domain.LambdaResponse, error) {
	objectKey, err := su.srtRepository.UploadFileToS3(request)
	if err != nil {
		return nil, err
	}

	request.FileName = objectKey

	response, err := su.srtRepository.TriggerLambdaFunc(request)
	if err != nil {
		return nil, err
	}

	userID, err := bson.ObjectIDFromHex("678f03edcc89ec934b05abf7")
	if err != nil {
		return nil, err
	}

	srtHistory := domain.SRTHistory{
		UserID:   userID,
		FileName: strings.Replace(request.FileHeader.Filename, ".mp4", ".srt", 1),
		S3URL:    response.Body.SRTURL,
		Duration: response.Body.Duration,
	}

	if err = su.CreateHistory(srtHistory); err != nil {
		return nil, err
	}

	return response, nil
}

func (su *srtUseCase) CreateHistory(srtHistory domain.SRTHistory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srtHistory.CreatedAt = time.Now().UTC()
	srtHistory.UpdatedAt = time.Now().UTC()
	if err := srtHistory.Validate(); err != nil {
		return err
	}

	return su.srtRepository.CreateHistory(ctx, srtHistory)
}
