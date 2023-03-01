package db

import (
	"context"
	"fmt"
	"os"
	"time"
	"xl-app/errs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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
		fmt.Println(o)
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

func (d *Db) CreateTable(tableName string) {
	_, tableErr := d.client.CreateTable(d.ctx, d.createTableInput(tableName))
	errs.PanicOnErr(tableErr)
	d.waitForTable(d.ctx, tableName)
	fmt.Println("table is ready...")
}

func (d *Db) waitForTable(ctx context.Context, tableName string) {
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
		os.Exit(1)
	}
}
