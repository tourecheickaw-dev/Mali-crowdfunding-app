package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// SetReply saves (or removes, if reply is empty) the organizer's reply to a donation comment.
// Targeted update: the status and the rest of the donation are never touched.
func (r *DonationRepository) SetReply(ctx context.Context, id, reply string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(#id)"),
		ExpressionAttributeNames: map[string]string{
			"#id": "id",
			"#r":  "reply",
			"#ra": "reply_at",
		},
	}

	if reply == "" {
		input.UpdateExpression = aws.String("REMOVE #r, #ra")
	} else {
		input.UpdateExpression = aws.String("SET #r = :r, #ra = :ra")
		input.ExpressionAttributeValues = map[string]types.AttributeValue{
			":r":  &types.AttributeValueMemberS{Value: reply},
			":ra": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339Nano)},
		}
	}

	if _, err := r.client.UpdateItem(ctx, input); err != nil {
		return fmt.Errorf("failed to save reply: %w", err)
	}
	return nil
}
