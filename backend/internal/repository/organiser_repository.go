package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

type OrganizerRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewOrganizerRepository(client *dynamodb.Client, tableName string) *OrganizerRepository {
	return &OrganizerRepository{
		client:    client,
		tableName: tableName,
	}
}

// SaveOrganizer creates or updates an organizer profile
func (r *OrganizerRepository) SaveOrganizer(ctx context.Context, organizer *models.Organizer) error {
	item, err := marshalItem(organizer)
	if err != nil {
		return fmt.Errorf("failed to marshal organizer: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save organizer: %w", err)
	}

	return nil
}

// GetOrganizerByID retrieves an organizer profile by unique ID
func (r *OrganizerRepository) GetOrganizerByID(ctx context.Context, id string) (*models.Organizer, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizer: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var organizer models.Organizer
	err = unmarshalItem(result.Item, &organizer)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal organizer: %w", err)
	}

	return &organizer, nil
}

// GetOrganizerByPhone retrieves an organizer by phone number using GSI 'phone-index'
func (r *OrganizerRepository) GetOrganizerByPhone(ctx context.Context, phone string) (*models.Organizer, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("phone-index"),
		KeyConditionExpression: aws.String("phone = :phone"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":phone": &types.AttributeValueMemberS{Value: phone},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query organizer by phone: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var organizer models.Organizer
	err = unmarshalItem(result.Items[0], &organizer)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal organizer: %w", err)
	}

	return &organizer, nil
}

// UpdateOrganizerStatus updates the verification status (KYC) of an organizer
func (r *OrganizerRepository) UpdateOrganizerStatus(ctx context.Context, id string, status models.OrganizerVerificationStatus) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET verification_status = :status, updated_at = :updated_at"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":     &types.AttributeValueMemberS{Value: string(status)},
			":updated_at": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	}

	_, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update organizer verification status: %w", err)
	}

	return nil
}

// AddToTotalRaised ATOMICALLY adds an amount to the lifetime total raised.
// The condition prevents accidentally creating an organizer that does not exist.
func (r *OrganizerRepository) AddToTotalRaised(ctx context.Context, id string, amount int64) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:    aws.String("ADD total_raised_lifetime :amt"),
		ConditionExpression: aws.String("attribute_exists(#id)"),
		ExpressionAttributeNames: map[string]string{
			"#id": "id",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":amt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", amount)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update organizer total: %w", err)
	}
	return nil
}
