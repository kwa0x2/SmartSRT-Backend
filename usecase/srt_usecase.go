package usecase

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

type srtUseCase struct {
	srtRepository     domain.SRTRepository
	usageUseCase      domain.UsageUseCase
	srtBaseRepository domain.BaseRepository[*domain.SRTHistory]
}

func NewSRTUseCase(srtRepository domain.SRTRepository, usageUseCase domain.UsageUseCase, srtBaseRepository domain.BaseRepository[*domain.SRTHistory]) domain.SRTUseCase {
	return &srtUseCase{
		srtRepository:     srtRepository,
		usageUseCase:      usageUseCase,
		srtBaseRepository: srtBaseRepository,
	}
}

func (su *srtUseCase) UploadFileAndConvertToSRT(request domain.FileConversionRequest) (*domain.LambdaResponse, error) {
	canUpload, err := su.usageUseCase.CheckUsageLimit(request.UserID, request.FileDuration)
	if err != nil {
		return nil, err
	}

	if !canUpload {
		return nil, utils.ErrLimitReached
	}

	objectKey, err := su.srtRepository.UploadFileToS3(request)
	if err != nil {
		return nil, err
	}

	request.FileName = objectKey

	response, err := su.srtRepository.TriggerLambdaFunc(request)
	if err != nil {
		return nil, err
	}

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := su.srtBaseRepository.GetDatabase().Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		if err = su.usageUseCase.UpdateUsage(txCtx, request.UserID, request.FileDuration); err != nil {
			return nil, err
		}

		fileType := filepath.Ext(request.FileHeader.Filename)
		srtHistory := &domain.SRTHistory{
			UserID:              request.UserID,
			FileName:            strings.Replace(request.FileHeader.Filename, fileType, ".srt", 1),
			S3URL:               response.Body.SRTURL,
			Duration:            request.FileDuration,
			WordsPerLine:        request.WordsPerLine,
			Punctuation:         request.Punctuation,
			ConsiderPunctuation: request.ConsiderPunctuation,
			CreatedAt:           time.Now().UTC(),
			UpdatedAt:           time.Now().UTC(),
		}

		if err = srtHistory.Validate(); err != nil {
			return nil, err
		}

		if err = su.srtBaseRepository.Create(txCtx, srtHistory); err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return nil, abortErr
		}
		return nil, err
	}

	return response, nil
}

func (su *srtUseCase) FindHistoriesByUserID(userID bson.ObjectID) ([]*domain.SRTHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	filter := bson.D{{Key: "user_id", Value: userID}}
	result, err := su.srtBaseRepository.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	return result, err
}
