package tagdb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TagRow struct {
	Tag       string                  `dynamodbav:"tag"`
	ObjectKey string                  `dynamodbav:"object_key"`
	UpdatedAt attributevalue.UnixTime `dynamodbav:"updated_at"`
}

func (r *TagRow) GetKey() map[string]types.AttributeValue {
	tag, err := attributevalue.Marshal(r.Tag)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"tag": tag}
}
