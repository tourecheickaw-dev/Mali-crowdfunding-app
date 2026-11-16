package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DB contient le client DynamoDB initialisé
type DB struct {
	Client *dynamodb.Client
}

// InitDynamoDB charge la configuration AWS et retourne le client
func InitDynamoDB(ctx context.Context, region string) (*DB, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("impossible de charger la config AWS: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DB{Client: client}, nil
}
