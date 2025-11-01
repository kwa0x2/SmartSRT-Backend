package main

import (
	"bytes"
	"log/slog"
	"mime/multipart"
	"os"

	"github.com/kwa0x2/SmartSRT-Backend/bootstrap"
	"github.com/kwa0x2/SmartSRT-Backend/config"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/rabbitmq"
	"github.com/kwa0x2/SmartSRT-Backend/repository"
	"github.com/kwa0x2/SmartSRT-Backend/usecase"
)

type fileReader struct {
	*bytes.Reader
}

func (f *fileReader) Close() error {
	return nil
}

type Consumer struct {
	env           *config.Env
	logger        *slog.Logger
	SRTUseCase    domain.SRTUseCase
	resendUseCase domain.ResendUseCase
	rabbitMQ      *domain.RabbitMQ
}

func NewConsumer(env *config.Env, logger *slog.Logger, SRTUseCase domain.SRTUseCase, ResendUseCase domain.ResendUseCase, rabbitMQ *domain.RabbitMQ) *Consumer {
	return &Consumer{
		env:           env,
		logger:        logger,
		SRTUseCase:    SRTUseCase,
		resendUseCase: ResendUseCase,
		rabbitMQ:      rabbitMQ,
	}
}

func (c *Consumer) Start() error {
	err := rabbitmq.StartWorkerPool(c.rabbitMQ, 5, func(msg domain.ConversionMessage) (*domain.LambdaResponse, error) {
		c.logger.Info("File conversion process started",
			slog.String("file_id", msg.FileID),
			slog.String("user_id", msg.UserID.Hex()),
			slog.String("file_name", msg.FileName),
			slog.Int64("file_size", msg.FileSize),
			slog.Float64("file_duration", msg.FileDuration),
		)

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

		if _, err := c.resendUseCase.SendSRTCreatedEmail(msg.Email, response.Body.SRTURL); err != nil {
			c.logger.Error("Email sending failed",
				slog.String("email", msg.Email),
				slog.String("file_id", msg.FileID),
				slog.String("error", err.Error()),
			)
		} else {
			c.logger.Info("Email sent successfully",
				slog.String("email", msg.Email),
				slog.String("file_id", msg.FileID),
				slog.String("srt_url", response.Body.SRTURL),
			)
		}

		c.logger.Info("File processed successfully",
			slog.String("file_id", msg.FileID),
			slog.String("user_id", msg.UserID.Hex()),
			slog.String("srt_url", response.Body.SRTURL),
		)
		return response, nil
	})

	if err != nil {
		c.logger.Error("Worker pool startup failed",
			slog.String("error", err.Error()),
		)
		return err
	}

	c.logger.Info("Consumer started successfully",
		slog.String("status", "waiting_for_messages"),
	)
	select {}
}

func main() {
	app := bootstrap.App()
	env := app.Env
	logger := slog.Default()

	db := app.MongoDatabase
	s3Client := app.S3Client
	lambdaClient := app.LambdaClient

	rabbitMQ, err := bootstrap.NewRabbitMQ(env)
	if err != nil {
		logger.Error("RabbitMQ connection failed",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer rabbitmq.Close(rabbitMQ)

	logger.Info("RabbitMQ connection established",
		slog.String("status", "connected"),
	)

	sr := repository.NewSRTRepository(s3Client, lambdaClient, db, env.AWSS3BucketName, env.AWSLambdaFuncName, domain.CollectionSRTHistory)
	usguc := usecase.NewUsageUseCase(env, repository.NewBaseRepository[*domain.Usage](db), repository.NewBaseRepository[*domain.User](db))
	srtUseCase := usecase.NewSRTUseCase(sr, usguc, repository.NewBaseRepository[*domain.SRTHistory](db))
	resendUseCase := usecase.NewResendUseCase(repository.NewResendRepository(app.ResendClient))

	consumer := NewConsumer(env, logger, srtUseCase, resendUseCase, rabbitMQ)
	if err = consumer.Start(); err != nil {
		logger.Error("Consumer error",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
