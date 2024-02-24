package tagdb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TagRow struct {
	Type      string                  `dynamodbav:"type"`       // PK
	Tag       string                  `dynamodbav:"tag"`        // SK, GSI1SK
	ObjectKey string                  `dynamodbav:"object_key"` // GSI1PK
	UpdatedAt attributevalue.UnixTime `dynamodbav:"updated_at"`
}

func (r *TagRow) GetKey() map[string]types.AttributeValue {
	tag, err := attributevalue.Marshal(r.Tag)
	if err != nil {
		panic(err)
	}
	typ, err := attributevalue.Marshal(r.Type)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{
		"type": typ,
		"tag":  tag,
	}
}
