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

type DonationRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDonationRepository(client *dynamodb.Client, tableName string) *DonationRepository {
	return &DonationRepository{
		client:    client,
		tableName: tableName,
	}
}

// CreateDonation inserts a new donation record into the database
func (r *DonationRepository) CreateDonation(ctx context.Context, donation *models.Donation) error {
	item, err := marshalItem(donation)
	if err != nil {
		return fmt.Errorf("failed to marshal donation: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save donation: %w", err)
	}

	return nil
}

// GetDonationByID retrieves a donation record by its unique ID
func (r *DonationRepository) GetDonationByID(ctx context.Context, id string) (*models.Donation, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get donation: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var donation models.Donation
	err = unmarshalItem(result.Item, &donation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal donation: %w", err)
	}

	return &donation, nil
}

// UpdateDonationStatus securely updates the status of a donation
func (r *DonationRepository) UpdateDonationStatus(ctx context.Context, id string, status models.DonationStatus) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET #s = :status, updated_at = :updated_at"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":     &types.AttributeValueMemberS{Value: string(status)},
			":updated_at": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	}

	_, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update donation status: %w", err)
	}

	return nil
}

// GetDonationsByCampaignID returns all donations of a campaign (campaign_id-index).
func (r *DonationRepository) GetDonationsByCampaignID(ctx context.Context, campaignID string) ([]models.Donation, error) {
	var donations []models.Donation
	var lastKey map[string]types.AttributeValue

	for {
		result, err := r.client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			IndexName:              aws.String("campaign_id-index"),
			KeyConditionExpression: aws.String("campaign_id = :cid"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":cid": &types.AttributeValueMemberS{Value: campaignID},
			},
			ExclusiveStartKey: lastKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query campaign donations: %w", err)
		}

		var page []models.Donation
		if err := unmarshalItems(result.Items, &page); err != nil {
			return nil, fmt.Errorf("failed to unmarshal donations: %w", err)
		}
		donations = append(donations, page...)

		if result.LastEvaluatedKey == nil {
			break
		}
		lastKey = result.LastEvaluatedKey
	}

	return donations, nil
}

// queryByIndex returns donations where the keyAttr attribute equals keyValue, using indexName.
func (r *DonationRepository) queryByIndex(ctx context.Context, indexName, keyAttr, keyValue string) ([]models.Donation, error) {
	var donations []models.Donation
	var lastKey map[string]types.AttributeValue
	for {
		result, err := r.client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			IndexName:              aws.String(indexName),
			KeyConditionExpression: aws.String(keyAttr + " = :v"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": &types.AttributeValueMemberS{Value: keyValue},
			},
			ExclusiveStartKey: lastKey,
		})
		if err != nil {
			return nil, fmt.Errorf("error while searching for donations (%s): %w", indexName, err)
		}
		var page []models.Donation
		if err := unmarshalItems(result.Items, &page); err != nil {
			return nil, fmt.Errorf("error unmarshaling donations: %w", err)
		}
		donations = append(donations, page...)
		if result.LastEvaluatedKey == nil {
			break
		}
		lastKey = result.LastEvaluatedKey
	}
	return donations, nil
}

// GetDonationsByDonorID returns all donations made by a donor (using donor_id-index).
func (r *DonationRepository) GetDonationsByDonorID(ctx context.Context, donorID string) ([]models.Donation, error) {
	return r.queryByIndex(ctx, "donor_id-index", "donor_id", donorID)
}

// UpdateDonorPhoto updates (or removes, if photoURL is empty) the photo displayed on a donation.
// Targeted update of a single field, to never overwrite the status or the rest of the donation.
func (r *DonationRepository) UpdateDonorPhoto(ctx context.Context, id, photoURL string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(#id)"),
		ExpressionAttributeNames: map[string]string{
			"#id": "id",
			"#p":  "donor_photo_url",
		},
	}

	if photoURL == "" {
		input.UpdateExpression = aws.String("REMOVE #p")
	} else {
		input.UpdateExpression = aws.String("SET #p = :p")
		input.ExpressionAttributeValues = map[string]types.AttributeValue{
			":p": &types.AttributeValueMemberS{Value: photoURL},
		}
	}

	if _, err := r.client.UpdateItem(ctx, input); err != nil {
		return fmt.Errorf("error updating donation photo: %w", err)
	}
	return nil
}

// UpdateComment updates (or removes, if comment is empty) a donation's comment.
// Without a comment, the organizer's reply is also removed.
func (r *DonationRepository) UpdateComment(ctx context.Context, id, comment string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(#id)"),
		ExpressionAttributeNames: map[string]string{
			"#id": "id",
			"#c":  "comment", // "comment" is a reserved word in DynamoDB
		},
	}

	if comment == "" {
		input.ExpressionAttributeNames["#r"] = "reply"
		input.ExpressionAttributeNames["#ra"] = "reply_at"
		input.UpdateExpression = aws.String("REMOVE #c, #r, #ra")
	} else {
		input.UpdateExpression = aws.String("SET #c = :c")
		input.ExpressionAttributeValues = map[string]types.AttributeValue{
			":c": &types.AttributeValueMemberS{Value: comment},
		}
	}

	if _, err := r.client.UpdateItem(ctx, input); err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

// SetCommentHidden hides or unhides a donation's comment.
func (r *DonationRepository) SetCommentHidden(ctx context.Context, id string, hidden bool) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(#id)"),
		ExpressionAttributeNames: map[string]string{
			"#id": "id",
			"#h":  "comment_hidden",
		},
	}

	if hidden {
		input.UpdateExpression = aws.String("SET #h = :h")
		input.ExpressionAttributeValues = map[string]types.AttributeValue{
			":h": &types.AttributeValueMemberBOOL{Value: true},
		}
	} else {
		input.UpdateExpression = aws.String("REMOVE #h")
	}

	if _, err := r.client.UpdateItem(ctx, input); err != nil {
		return fmt.Errorf("failed to change comment visibility: %w", err)
	}
	return nil
}
