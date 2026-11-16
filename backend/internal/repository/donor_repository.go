package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

type DonorRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDonorRepository(client *dynamodb.Client, tableName string) *DonorRepository {
	return &DonorRepository{
		client:    client,
		tableName: tableName,
	}
}

// SaveDonor creates or updates a donor record in the DynamoDB table
func (r *DonorRepository) SaveDonor(ctx context.Context, donor *models.Donor) error {
	item, err := marshalItem(donor)
	if err != nil {
		return fmt.Errorf("failed to marshal donor: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save donor: %w", err)
	}

	return nil
}

// GetDonorByID retrieves a donor record by unique ID
func (r *DonorRepository) GetDonorByID(ctx context.Context, id string) (*models.Donor, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get donor: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var donor models.Donor
	err = unmarshalItem(result.Item, &donor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal donor: %w", err)
	}

	return &donor, nil
}

// GetDonorByPhone retrieves a donor record by phone number using GSI 'phone-index'
func (r *DonorRepository) GetDonorByPhone(ctx context.Context, phone string) (*models.Donor, error) {
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
		return nil, fmt.Errorf("failed to query donor by phone: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var donor models.Donor
	err = unmarshalItem(result.Items[0], &donor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal donor: %w", err)
	}

	return &donor, nil
}
