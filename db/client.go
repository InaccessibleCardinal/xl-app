package db

import (
	"context"
	"strconv"

	"fmt"
	"time"
	"xl-app/errs"
	T "xl-app/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Db struct {
	ctx    context.Context
	client *dynamodb.Client
}

func New(ctx context.Context) *Db {
	client := createDynamoClient(ctx)
	return &Db{ctx: ctx, client: client}
}

func createAwsConfig(ctx context.Context) aws.Config {
	cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	errs.PanicOnErr(err)
	return cfg
}

func createDynamoClient(ctx context.Context) *dynamodb.Client {
	cfg := createAwsConfig(ctx)
	return dynamodb.NewFromConfig(cfg)
}

func (d *Db) createTableInput(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
	}
}

func (d *Db) CreateTable(tableName string) error {
	_, tableErr := d.client.CreateTable(d.ctx, d.createTableInput(tableName))
	errs.PanicOnErr(tableErr)
	if err := d.waitForTable(tableName); err != nil {
		return err
	}
	fmt.Println("table is ready...")
	return nil
}

func (d *Db) waitForTable(tableName string) error {
	w := dynamodb.NewTableExistsWaiter(d.client)
	err := w.Wait(d.ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		},
		2*time.Minute,
		func(o *dynamodb.TableExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		fmt.Printf("Waiter timed out waiting for table %s, error: %s", tableName, err.Error())
		return err
	}
	return nil
}

func addIds(xlData []T.StringMap) []T.StringMap {
	var val []map[string]string
	for i, item := range xlData {
		newMap := map[string]string{"id": strconv.Itoa(i)}
		for k, v := range item {
			newMap[k] = v
		}
		val = append(val, newMap)
	}
	return val
}

func (d *Db) BulkSave(xlDto T.XLDto) error {
	batch := make(map[string][]types.WriteRequest)
	itemsWithIds := addIds(xlDto.XlData)
	var requests []types.WriteRequest
	for _, item := range itemsWithIds {
		marshaledItem, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}
		requests = append(requests, types.WriteRequest{PutRequest: &types.PutRequest{Item: marshaledItem}})
	}
	batch[xlDto.DbName] = requests
	out, err := d.client.BatchWriteItem(d.ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: batch,
	})
	if err != nil {
		fmt.Printf("error writing %s", err.Error())
		return err
	}
	if len(out.UnprocessedItems) != 0 {
		fmt.Println("there were ", len(out.UnprocessedItems), " unprocessed records")
	}
	return nil
}
