package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(9000) + 1000

	newFileName := fmt.Sprintf("%s_%s_%s", "autosrt.com", strconv.Itoa(randomNumber), request.FileHeader.Filename)
	objectKey := fmt.Sprintf("videos/%s/%s", request.UserID.Hex(), newFileName)

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

func (sr *srtRepository) CreateHistory(ctx context.Context, srtHistory domain.SRTHistory) error {
	result, err := sr.collection.InsertOne(ctx, srtHistory)
	if err != nil {
		return err
	}

	srtHistory.ID = result.InsertedID.(bson.ObjectID)

	return nil
}

func (sr *srtRepository) FindHistories(ctx context.Context, filter bson.D, opts *options.FindOptionsBuilder) ([]domain.SRTHistory, error) {
	cursor, err := sr.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var srtHistories []domain.SRTHistory

	if err = cursor.All(ctx, &srtHistories); err != nil {
		return nil, err
	}

	return srtHistories, nil
}
