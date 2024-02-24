package tagdb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Client struct {
	dynamodbCli *dynamodb.Client
	tableName   string
}

func NewClient(cfg aws.Config, tableName string) *Client {
	dynamoDBCli := dynamodb.NewFromConfig(cfg)
	return &Client{
		dynamodbCli: dynamoDBCli,
		tableName:   tableName,
	}
}

type PutTagsInput struct {
	Tags      []string
	ObjectKey string
	UpdatedAt time.Time
}

func (c *Client) PutTags(ctx context.Context, input *PutTagsInput) error {
	for _, tag := range input.Tags {
		tagRow := TagRow{
			Tag:       tag,
			ObjectKey: input.ObjectKey,
			UpdatedAt: attributevalue.UnixTime(input.UpdatedAt),
		}
		tagItem, err := attributevalue.MarshalMap(tagRow)
		if err != nil {
			return err
		}

		if _, err := c.dynamodbCli.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(c.tableName),
			Item:      tagItem,
		}); err != nil {
			return err
		}
	}
	return nil
}

type FindByTagInput struct {
	Tag string
}

func (c *Client) FindByTag(ctx context.Context, input *FindByTagInput) (*TagRow, error) {
	getItemResult, err := c.dynamodbCli.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]types.AttributeValue{
			"tag": &types.AttributeValueMemberS{
				Value: input.Tag,
			},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	tagRow := TagRow{}
	if err = attributevalue.UnmarshalMap(getItemResult.Item, &tagRow); err != nil {
		return nil, err
	}
	return &tagRow, nil
}

func (c *Client) GetAll(ctx context.Context) ([]*TagRow, error) {
	tagRows := []*TagRow{}
	scanPaginator := dynamodb.NewScanPaginator(c.dynamodbCli, &dynamodb.ScanInput{
		TableName:      aws.String(c.tableName),
		ConsistentRead: aws.Bool(true),
	})
	for scanPaginator.HasMorePages() {
		response, err := scanPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		page := []*TagRow{}
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &page); err != nil {
			return nil, err
		}
		tagRows = append(tagRows, page...)
	}
	return tagRows, nil
}
