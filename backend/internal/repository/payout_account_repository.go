package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

type PayoutAccountRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewPayoutAccountRepository(client *dynamodb.Client, tableName string) *PayoutAccountRepository {
	return &PayoutAccountRepository{
		client:    client,
		tableName: tableName,
	}
}

// SavePayoutAccount creates or updates a payout account
func (r *PayoutAccountRepository) SavePayoutAccount(ctx context.Context, account *models.PayoutAccount) error {
	item, err := marshalItem(account)
	if err != nil {
		return fmt.Errorf("failed to marshal payout account: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save payout account: %w", err)
	}

	return nil
}

// GetPayoutAccountByID retrieves a payout account by its unique ID
func (r *PayoutAccountRepository) GetPayoutAccountByID(ctx context.Context, id string) (*models.PayoutAccount, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get payout account: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var account models.PayoutAccount
	err = unmarshalItem(result.Item, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payout account: %w", err)
	}

	return &account, nil
}

// GetPayoutAccountsByOrganizerID retrieves all payout accounts belonging to an organizer
func (r *PayoutAccountRepository) GetPayoutAccountsByOrganizerID(ctx context.Context, organizerID string) ([]models.PayoutAccount, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("organizer_id-index"),
		KeyConditionExpression: aws.String("organizer_id = :organizer_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":organizer_id": &types.AttributeValueMemberS{Value: organizerID},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query payout accounts by organizer ID: %w", err)
	}

	var accounts []models.PayoutAccount
	err = unmarshalItems(result.Items, &accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payout accounts list: %w", err)
	}

	return accounts, nil
}

// SetActivePayoutAccount activates a target payout account and deactivates all other accounts for the organizer
func (r *PayoutAccountRepository) SetActivePayoutAccount(ctx context.Context, organizerID, accountIDToActivate string) error {
	accounts, err := r.GetPayoutAccountsByOrganizerID(ctx, organizerID)
	if err != nil {
		return fmt.Errorf("failed to fetch accounts before activation: %w", err)
	}

	for _, acc := range accounts {
		shouldBeActive := acc.ID == accountIDToActivate

		_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(r.tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: acc.ID},
			},
			UpdateExpression: aws.String("SET is_active = :active"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":active": &types.AttributeValueMemberBOOL{Value: shouldBeActive},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to update active status for account %s: %w", acc.ID, err)
		}
	}

	return nil
}
