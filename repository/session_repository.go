package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"log"
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

func (sr *sessionRepository) CreateSession(ctx context.Context, session domain.Session) error {
	av, err := attributevalue.MarshalMap(session)
	if err != nil {
		log.Fatalf("failed to marshal session: %v", err)
	}

	_, err = sr.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &sr.tableName,
		Item:      av,
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

	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, utils.ErrSessionNotFound
	}

	var session domain.Session
	err = attributevalue.UnmarshalMap(resp.Item, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil

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

func (sr *sessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := sr.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(sr.tableName),
		Key: map[string]types.AttributeValue{
			"session_id": &types.AttributeValueMemberS{Value: sessionID},
		},
	})

	return err
}
