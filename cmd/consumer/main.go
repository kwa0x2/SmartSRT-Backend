package main

import (
	"bytes"
	"log"
	"mime/multipart"

	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/rabbitmq"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
)

type fileReader struct {
	*bytes.Reader
}

func (f *fileReader) Close() error {
	return nil
}

type Consumer struct {
	env           *config.Env
	SRTUseCase    domain.SRTUseCase
	resendUseCase domain.ResendUseCase
	rabbitMQ      *domain.RabbitMQ
}

func NewConsumer(env *config.Env, SRTUseCase domain.SRTUseCase, ResendUseCase domain.ResendUseCase, rabbitMQ *domain.RabbitMQ) *Consumer {
	return &Consumer{
		env:           env,
		SRTUseCase:    SRTUseCase,
		resendUseCase: ResendUseCase,
		rabbitMQ:      rabbitMQ,
	}
}

func (c *Consumer) Start() error {
	log.Printf("Starting SRT conversion consumer...")

	err := rabbitmq.StartWorkerPool(c.rabbitMQ, 5, func(msg domain.ConversionMessage) (*domain.LambdaResponse, error) {
		log.Printf("Processing conversion for file %s by user %s", msg.FileID, msg.UserID)

		request := domain.FileConversionRequest{
			UserID:              msg.UserID,
			WordsPerLine:        msg.WordsPerLine,
			Punctuation:         msg.Punctuation,
			ConsiderPunctuation: msg.ConsiderPunctuation,
			FileName:            msg.FileName,
			File:                &fileReader{bytes.NewReader(msg.FileContent)},
			FileHeader: multipart.FileHeader{
				Filename: msg.FileName,
				Size:     msg.FileSize,
			},
			FileDuration: msg.FileDuration,
		}

		response, err := c.SRTUseCase.UploadFileAndConvertToSRT(request)
		if err != nil {
			return nil, err
		}

		go func() {
			if _, err := c.resendUseCase.SendSRTCreatedEmail(msg.Email, response.Body.SRTURL); err != nil {
				log.Printf("Error sending email: %v", err)
			} else {
				log.Printf("Email sent successfully to %s", msg.Email)
			}
		}()

		return response, nil
	})

	if err != nil {
		return err
	}

	log.Printf("Consumer started successfully, waiting for messages...")
	select {}
}

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.MongoDatabase
	s3Client := app.S3Client
	lambdaClient := app.LambdaClient

	rabbitMQ, err := bootstrap.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close(rabbitMQ)
	log.Printf("RabbitMQ connection established")

	sr := repository.NewSRTRepository(s3Client, lambdaClient, db, env.AWSS3BucketName, env.AWSLambdaFuncName, domain.CollectionSRTHistory)
	usguc := usecase.NewUsageUseCase(repository.NewBaseRepository[*domain.Usage](db), repository.NewBaseRepository[*domain.User](db))
	srtUseCase := usecase.NewSRTUseCase(sr, usguc, repository.NewBaseRepository[*domain.SRTHistory](db))
	resendUseCase := usecase.NewResendUseCase(repository.NewResendRepository(app.ResendClient))

	consumer := NewConsumer(env, srtUseCase, resendUseCase, rabbitMQ)
	if err = consumer.Start(); err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
