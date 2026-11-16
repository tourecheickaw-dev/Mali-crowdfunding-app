package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

type MediaRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewMediaRepository(client *dynamodb.Client, tableName string) *MediaRepository {
	return &MediaRepository{
		client:    client,
		tableName: tableName,
	}
}

// SaveMedia inserts or updates media file metadata in DynamoDB
func (r *MediaRepository) SaveMedia(ctx context.Context, media *models.CampaignMedia) error {
	item, err := marshalItem(media)
	if err != nil {
		return fmt.Errorf("failed to marshal media: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save media: %w", err)
	}

	return nil
}

// GetMediaByCampaignID retrieves all media items linked to a campaign (via 'campaign_id-index' GSI)
func (r *MediaRepository) GetMediaByCampaignID(ctx context.Context, campaignID string) ([]models.CampaignMedia, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("campaign_id-index"),
		KeyConditionExpression: aws.String("campaign_id = :campaign_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":campaign_id": &types.AttributeValueMemberS{Value: campaignID},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query media by campaign ID: %w", err)
	}

	var mediaList []models.CampaignMedia
	err = unmarshalItems(result.Items, &mediaList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal media list: %w", err)
	}

	return mediaList, nil
}

// DeleteMedia deletes a specific media record by its unique ID
func (r *MediaRepository) DeleteMedia(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := r.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}

	return nil
}
