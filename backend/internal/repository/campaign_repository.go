package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

type CampaignRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewCampaignRepository(client *dynamodb.Client, tableName string) *CampaignRepository {
	return &CampaignRepository{
		client:    client,
		tableName: tableName,
	}
}

// CreateCampaign inserts a campaign record into DynamoDB
func (r *CampaignRepository) CreateCampaign(ctx context.Context, campaign models.Campaign) error {
	item, err := marshalItem(campaign)
	if err != nil {
		return fmt.Errorf("failed to marshal campaign: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to insert campaign into DynamoDB: %w", err)
	}

	return nil
}

// GetCampaignByID retrieves a campaign record by its ID
func (r *CampaignRepository) GetCampaignByID(ctx context.Context, id string) (*models.Campaign, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var campaign models.Campaign
	if err := unmarshalItem(result.Item, &campaign); err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaign: %w", err)
	}

	return &campaign, nil
}

// GetCampaignsByOrganizer retrieves all campaigns belonging to an organizer
// Uses the "organizer_id-index" GSI created on the table
func (r *CampaignRepository) GetCampaignsByOrganizer(ctx context.Context, organizerID string) ([]models.Campaign, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("organizer_id-index"),
		KeyConditionExpression: aws.String("organizer_id = :oid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":oid": &types.AttributeValueMemberS{Value: organizerID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query campaigns by organizer: %w", err)
	}

	var campaigns []models.Campaign
	if err := unmarshalItems(result.Items, &campaigns); err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaigns list: %w", err)
	}

	return campaigns, nil
}

// GetCampaignBySlug retrieves a campaign record by its URL slug
// Uses the "url_slug-index" GSI created on the table
func (r *CampaignRepository) GetCampaignBySlug(ctx context.Context, slug string) (*models.Campaign, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("url_slug-index"),
		KeyConditionExpression: aws.String("url_slug = :slug"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":slug": &types.AttributeValueMemberS{Value: slug},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query campaign by slug: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var campaigns []models.Campaign
	if err := unmarshalItems(result.Items, &campaigns); err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaigns list: %w", err)
	}

	return &campaigns[0], nil
}

// UpdateCampaignAmount ATOMICALLY increments current_amount.
// Never use PutItem here: two concurrent donations would overwrite each other.
func (r *CampaignRepository) UpdateCampaignAmount(ctx context.Context, campaignID string, amountToAdd int64) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: campaignID},
		},
		UpdateExpression: aws.String("ADD current_amount :amt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":amt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", amountToAdd)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update campaign amount: %w", err)
	}

	return nil
}

// ListActiveCampaigns returns active campaigns (at most max).
// A Scan is fine for the MVP (a few hundred campaigns). Replace it with a status index later.
func (r *CampaignRepository) ListActiveCampaigns(ctx context.Context, max int) ([]models.Campaign, error) {
	var out []models.Campaign
	var startKey map[string]types.AttributeValue

	for {
		res, err := r.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:                aws.String(r.tableName),
			FilterExpression:         aws.String("#st = :active"),
			ExpressionAttributeNames: map[string]string{"#st": "status"},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":active": &types.AttributeValueMemberS{Value: string(models.CampaignActive)},
			},
			ExclusiveStartKey: startKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaigns: %w", err)
		}

		var page []models.Campaign
		if err := unmarshalItems(res.Items, &page); err != nil {
			return nil, fmt.Errorf("failed to unmarshal campaigns page: %w", err)
		}
		out = append(out, page...)

		if res.LastEvaluatedKey == nil || len(out) >= max {
			break
		}
		startKey = res.LastEvaluatedKey
	}
	return out, nil
}
