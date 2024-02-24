package tagdb

import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

type TagRow struct {
	Tag       string                  `dynamodbav:"tag"`
	ObjectKey string                  `dynamodbav:"object_key"`
	UpdatedAt attributevalue.UnixTime `dynamodbav:"updated_at"`
}
