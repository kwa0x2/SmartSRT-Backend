package usecase

import (
	"github.com/kwa0x2/AutoSRT-Backend/domain"
)

type subtitleUseCase struct {
	subtitleRepository domain.SubtitleRepository
}

func NewSubtitleUseCase(subtitleRepository domain.SubtitleRepository) domain.SubtitleUseCase {
	return &subtitleUseCase{
		subtitleRepository: subtitleRepository,
	}
}

func (su *subtitleUseCase) UploadFileAndConvertToSRT(request domain.FileConversionRequest) (*domain.LambdaResponse, error) {
	objectKey, err := su.subtitleRepository.UploadFileToS3(request)
	if err != nil {
		return nil, err
	}

	request.FileName = objectKey

	response, err := su.subtitleRepository.TriggerLambdaFunc(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
