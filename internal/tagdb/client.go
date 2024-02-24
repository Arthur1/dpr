package tagdb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Client interface {
	PutMultiple(ctx context.Context, rows []*TagRow) error
	FindByTag(ctx context.Context, tag string) (*TagRow, error)
	FindByDigest(ctx context.Context, digest string) (*TagRow, error)
	GetAll(ctx context.Context) ([]*TagRow, error)
	GetByTagPrefix(ctx context.Context, tagPrefix string) ([]*TagRow, error)
	GetByObjectKey(ctx context.Context, objectKey string) ([]*TagRow, error)
	DeleteMultiple(ctx context.Context, rows []*TagRow) error
}

type ClientImpl struct {
	dynamodbCli *dynamodb.Client
	tableName   string
}

var _ Client = new(ClientImpl)

func NewClientImpl(cfg aws.Config, tableName string) *ClientImpl {
	dynamoDBCli := dynamodb.NewFromConfig(cfg)
	return &ClientImpl{
		dynamodbCli: dynamoDBCli,
		tableName:   tableName,
	}
}

type PutMultipleParam struct {
	Tags      []string
	ObjectKey string
	UpdatedAt time.Time
}

func (c *ClientImpl) PutMultiple(ctx context.Context, rows []*TagRow) error {
	for _, r := range rows {
		item, err := attributevalue.MarshalMap(r.ToDynamoDBTagRow())
		if err != nil {
			return err
		}
		if _, err := c.dynamodbCli.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(c.tableName),
			Item:      item,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientImpl) FindByTag(ctx context.Context, tag string) (*TagRow, error) {
	getItemResult, err := c.dynamodbCli.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]types.AttributeValue{
			"type": &types.AttributeValueMemberS{
				Value: "tag",
			},
			"tag": &types.AttributeValueMemberS{
				Value: tag,
			},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	drow := DynamoDBTagRow{}
	if err = attributevalue.UnmarshalMap(getItemResult.Item, &drow); err != nil {
		return nil, err
	}
	return drow.ToTagRow(), nil
}

func (c *ClientImpl) FindByDigest(ctx context.Context, digest string) (*TagRow, error) {
	getItemResult, err := c.dynamodbCli.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]types.AttributeValue{
			"type": &types.AttributeValueMemberS{
				Value: "digest",
			},
			"tag": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("@%s", digest),
			},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	drow := DynamoDBTagRow{}
	if err = attributevalue.UnmarshalMap(getItemResult.Item, &drow); err != nil {
		return nil, err
	}
	return drow.ToTagRow(), nil
}

func (c *ClientImpl) GetAll(ctx context.Context) ([]*TagRow, error) {
	rows := []*TagRow{}
	scanPaginator := dynamodb.NewScanPaginator(c.dynamodbCli, &dynamodb.ScanInput{
		TableName:      aws.String(c.tableName),
		ConsistentRead: aws.Bool(true),
	})
	for scanPaginator.HasMorePages() {
		response, err := scanPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		page := []*DynamoDBTagRow{}
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &page); err != nil {
			return nil, err
		}
		for _, drow := range page {
			rows = append(rows, drow.ToTagRow())
		}
	}
	return rows, nil
}

func (c *ClientImpl) GetByTagPrefix(ctx context.Context, tagPrefix string) ([]*TagRow, error) {
	keyEx := expression.Key("type").Equal(expression.Value("tag")).
		And(expression.Key("tag").BeginsWith(tagPrefix))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}

	rows := []*TagRow{}
	queryPaginator := dynamodb.NewQueryPaginator(c.dynamodbCli, &dynamodb.QueryInput{
		TableName:                 aws.String(c.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	for queryPaginator.HasMorePages() {
		response, err := queryPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		page := []*DynamoDBTagRow{}
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &page); err != nil {
			return nil, err
		}
		for _, drow := range page {
			rows = append(rows, drow.ToTagRow())
		}
	}
	return rows, nil
}

func (c *ClientImpl) GetByObjectKey(ctx context.Context, objectKey string) ([]*TagRow, error) {
	keyEx := expression.Key("object_key").Equal(expression.Value(objectKey))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}

	rows := []*TagRow{}
	queryPaginator := dynamodb.NewQueryPaginator(c.dynamodbCli, &dynamodb.QueryInput{
		TableName:                 aws.String(c.tableName),
		IndexName:                 aws.String("index_object_key_tag"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	for queryPaginator.HasMorePages() {
		response, err := queryPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		page := []*DynamoDBTagRow{}
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &page); err != nil {
			return nil, err
		}
		for _, drow := range page {
			rows = append(rows, drow.ToTagRow())
		}
	}
	return rows, nil
}

func (c *ClientImpl) DeleteMultiple(ctx context.Context, rows []*TagRow) error {
	for _, row := range rows {
		if _, err := c.dynamodbCli.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(c.tableName),
			Key:       row.ToDynamoDBTagRow().GetKey(),
		}); err != nil {
			return err
		}
	}
	return nil
}
