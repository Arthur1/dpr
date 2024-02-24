package tagdb

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TagRow struct {
	Type      string
	Tag       string
	ObjectKey string
	UpdatedAt time.Time
}

func (r *TagRow) ToDynamoDBTagRow() *DynamoDBTagRow {
	return &DynamoDBTagRow{
		Type:      r.Type,
		Tag:       r.Tag,
		ObjectKey: r.ObjectKey,
		UpdatedAt: attributevalue.UnixTime(r.UpdatedAt),
	}
}

type DynamoDBTagRow struct {
	Type      string                  `dynamodbav:"type"`       // PK
	Tag       string                  `dynamodbav:"tag"`        // SK, GSI1SK
	ObjectKey string                  `dynamodbav:"object_key"` // GSI1PK
	UpdatedAt attributevalue.UnixTime `dynamodbav:"updated_at"`
}

func (r *DynamoDBTagRow) ToTagRow() *TagRow {
	return &TagRow{
		Type:      r.Type,
		Tag:       r.Tag,
		ObjectKey: r.ObjectKey,
		UpdatedAt: time.Time(r.UpdatedAt),
	}
}

func (r *DynamoDBTagRow) GetKey() map[string]types.AttributeValue {
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
