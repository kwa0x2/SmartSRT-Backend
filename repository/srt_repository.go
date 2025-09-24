package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type srtRepository struct {
	s3Client       *s3.Client
	lambdaClient   *lambda.Client
	lambdaFuncName string
	bucketName     string
	collection     *mongo.Collection
}

func NewSRTRepository(s3Client *s3.Client, lambdaClient *lambda.Client, db *mongo.Database, bucketName, lambdaFuncName, collection string) domain.SRTRepository {
	return &srtRepository{
		s3Client:       s3Client,
		lambdaClient:   lambdaClient,
		lambdaFuncName: lambdaFuncName,
		bucketName:     bucketName,
		collection:     db.Collection(collection),
	}
}

func (sr *srtRepository) UploadFileToS3(request domain.FileConversionRequest) (string, error) {
	newFileName := fmt.Sprintf("%s_%d_%s", "smartsrt.com", time.Now().UTC().Unix(), request.FileHeader.Filename)
	objectKey := fmt.Sprintf("files/%s/%s", request.UserID.Hex(), newFileName)

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

func (sr *srtRepository) TriggerLambdaFunc(request domain.FileConversionRequest) (*domain.LambdaResponse, error) {
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
		return nil, fmt.Errorf("lambda function error: %s", *result.FunctionError)
	}

	var rawResponse domain.LambdaResponse

	if err = json.Unmarshal(result.Payload, &rawResponse); err != nil {
		return nil, err
	}

	if rawResponse.StatusCode != http.StatusOK {
		return nil, errors.New(rawResponse.Body.Message)
	}

	return &rawResponse, nil
}
