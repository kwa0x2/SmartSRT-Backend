package domain

import (
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	UserID              bson.ObjectID `json:"user_id"`
	WordsPerLine        int           `json:"words_per_line"`
	Punctuation         bool          `json:"punctuation"`
	ConsiderPunctuation bool          `json:"consider_punctuation"`
	FileName            string        `json:"file_name"`
	File                multipart.File
	FileHeader          multipart.FileHeader
	FileDuration        float64
}

const (
	CollectionSRTHistory = "srt_history"
)

type SRTHistory struct {
	ID                  bson.ObjectID `bson:"_id,omitempty"`
	UserID              bson.ObjectID `bson:"user_id" validate:"required"`
	FileName            string        `bson:"file_name" validate:"required"`
	S3URL               string        `bson:"s3_url" validate:"required"`
	Duration            float64       `bson:"duration"`
	WordsPerLine        int           `bson:"words_per_line"`
	Punctuation         bool          `bson:"punctuation"`
	ConsiderPunctuation bool          `bson:"consider_punctuation"`
	CreatedAt           time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt           time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt           *time.Time    `bson:"deleted_at,omitempty"`
}

func (s *SRTHistory) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (s *SRTHistory) GetCollectionName() string {
	return CollectionSRTHistory
}

func (s *SRTHistory) SetID(id bson.ObjectID) {
	s.ID = id
}

type SRTUseCase interface {
	UploadFileAndConvertToSRT(request FileConversionRequest) (*LambdaResponse, error)
	FindHistoriesByUserID(userID bson.ObjectID) ([]*SRTHistory, error)
}

type SRTRepository interface {
	UploadFileToS3(request FileConversionRequest) (string, error)
	TriggerLambdaFunc(request FileConversionRequest) (*LambdaResponse, error)
}
