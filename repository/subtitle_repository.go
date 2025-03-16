package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
)

type subtitleRepository struct {
	s3Client       *s3.Client
	lambdaClient   *lambda.Client
	lambdaFuncName string
	bucketName     string
}

func NewSubtitleRepository(s3Client *s3.Client, lambdaClient *lambda.Client, bucketName, lambdaFuncName string) domain.SubtitleRepository {
	return &subtitleRepository{
		s3Client:       s3Client,
		lambdaClient:   lambdaClient,
		lambdaFuncName: lambdaFuncName,
		bucketName:     bucketName,
	}
}

func (sr *subtitleRepository) UploadFileToS3(request domain.FileConversionRequest) (string, error) {
	newFileName := fmt.Sprintf("%s_%s", uuid.New().String(), request.FileHeader.Filename)
	objectKey := fmt.Sprintf("videos/%s/%s", request.UserID, newFileName)

	input := &s3.PutObjectInput{
		Bucket: aws.String(sr.bucketName),
		Key:    aws.String(objectKey),
		Body:   request.File,
	}

	_, err := sr.s3Client.PutObject(context.Background(), input)
	if err != nil {
		return "", err
	}

	return newFileName, nil
}

func (sr *subtitleRepository) TriggerLambdaFunc(request domain.FileConversionRequest) (*domain.LambdaResponse, error) {
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	input := &lambda.InvokeInput{
		FunctionName: aws.String(sr.lambdaFuncName),
		Payload:      jsonPayload,
	}

	result, err := sr.lambdaClient.Invoke(context.Background(), input)
	if err != nil {
		return nil, err
	}

	if result.FunctionError != nil {
		return nil, err
	}

	var rawResponse domain.LambdaResponse

	if err = json.Unmarshal(result.Payload, &rawResponse); err != nil {
		return nil, err
	}

	if rawResponse.StatusCode != 200 {
		return nil, err
	}

	return &rawResponse, nil
}
