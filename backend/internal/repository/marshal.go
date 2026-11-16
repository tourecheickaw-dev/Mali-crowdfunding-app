package repository

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Our models use `json` tags (id, organizer_id...). By default, the DynamoDB SDK only reads
// `dynamodbav` tags: without this option, fields would be stored as "ID",
// "OrganizerID"... and our table keys and indexes would be untraceable.

func marshalItem(in interface{}) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMapWithOptions(in, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
}

func unmarshalItem(item map[string]types.AttributeValue, out interface{}) error {
	return attributevalue.UnmarshalMapWithOptions(item, out, func(o *attributevalue.DecoderOptions) {
		o.TagKey = "json"
	})
}

func unmarshalItems(items []map[string]types.AttributeValue, out interface{}) error {
	return attributevalue.UnmarshalListOfMapsWithOptions(items, out, func(o *attributevalue.DecoderOptions) {
		o.TagKey = "json"
	})
}
