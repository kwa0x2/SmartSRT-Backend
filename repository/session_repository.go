package repository

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"strconv"
)

type sessionRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewSessionRepository(client *dynamodb.Client, tableName string) domain.SessionRepository {
	return &sessionRepository{
		client:    client,
		tableName: tableName,
	}
}

func (sr *sessionRepository) CreateSession(ctx context.Context, sessionID string, TTL int) error {
	_, err := sr.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &sr.tableName,
		Item: map[string]types.AttributeValue{
			"session_id": &types.AttributeValueMemberS{Value: sessionID},
			"ttl":        &types.AttributeValueMemberN{Value: strconv.Itoa(TTL)},
		},
	})
	return err
}

func (sr *sessionRepository) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	resp, err := sr.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(sr.tableName),
		Key: map[string]types.AttributeValue{
			"session_id": &types.AttributeValueMemberS{Value: sessionID},
		},
	})

	if err != nil || resp.Item == nil {
		return nil, err
	}

	TTLUnixStr, ok := resp.Item["ttl"].(*types.AttributeValueMemberN)
	if !ok {
		return nil, errors.New("ttl field is missing or invalid")
	}

	TTLUnix, atoiErr := strconv.Atoi(TTLUnixStr.Value)
	if atoiErr != nil {
		return nil, atoiErr
	}

	return &domain.Session{
		SessionID: sessionID,
		TTL:       TTLUnix,
	}, nil
}

func (sr *sessionRepository) UpdateSessionTTL(ctx context.Context, sessionID string, newTTL int) error {
	_, err := sr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(sr.tableName),
		Key: map[string]types.AttributeValue{
			"session_id": &types.AttributeValueMemberS{Value: sessionID},
		},
		UpdateExpression: aws.String("SET #ttl = :ttl"),
		ExpressionAttributeNames: map[string]string{
			"#ttl": "ttl",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ttl": &types.AttributeValueMemberN{Value: strconv.Itoa(newTTL)},
		},
		ReturnValues: types.ReturnValueAllNew,
	})

	if err != nil {
		return err
	}

	return nil
}
