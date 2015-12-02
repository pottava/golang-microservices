package aws

/**
 * @see https://github.com/aws/aws-sdk-go/blob/master/service/dynamodb/api.go
 */

import (
	"sync"

	appcfg "github.com/pottava/golang-micro-services/app-dbio/app/config"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamoCfg *awssdk.Config
var dynamoOnce sync.Once

func init() {
	dynamoOnce.Do(func() {
		dynamoCfg = config()
		if host := appcfg.NewConfig().DynamoDbLocal; host != "" {
			dynamoCfg.Endpoint = awssdk.String("http://" + host + ":8000")
		}
	})
}

// DynamoTables responses dynamodb tables
func DynamoTables() (clusters []*string, e error) {
	resp, err := dynamodb.New(session.New(), dynamoCfg).ListTables(nil)
	if err != nil {
		return nil, err
	}
	return resp.TableNames, nil
}

// DynamoTable responses a specific dynamodb table
func DynamoTable(name string) (table *dynamodb.TableDescription, e error) {
	resp, err := dynamodb.New(session.New(), dynamoCfg).DescribeTable(&dynamodb.DescribeTableInput{
		TableName: awssdk.String(name),
	})
	if err != nil {
		return nil, err
	}
	return resp.Table, nil
}

// DynamoScan responses dynamodb table records
func DynamoScan(name string) (records []map[string]*dynamodb.AttributeValue, count int64, e error) {
	resp, err := dynamodb.New(session.New(), dynamoCfg).Scan(&dynamodb.ScanInput{
		TableName: awssdk.String(name),
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.Items, *resp.ScannedCount, nil
}

// DynamoRecords responses a dynamodb specific records
func DynamoRecords(name, expression string, attributes map[string]*dynamodb.AttributeValue) (records []map[string]*dynamodb.AttributeValue, count int64, e error) {
	resp, err := dynamodb.New(session.New(), dynamoCfg).Query(&dynamodb.QueryInput{
		TableName:                 awssdk.String(name),
		KeyConditionExpression:    awssdk.String(expression),
		ExpressionAttributeValues: attributes,
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.Items, *resp.ScannedCount, nil
}

// DynamoRecord responses a dynamodb specific record
func DynamoRecord(name string, key map[string]*dynamodb.AttributeValue) (record map[string]*dynamodb.AttributeValue, e error) {
	resp, err := dynamodb.New(session.New(), dynamoCfg).GetItem(&dynamodb.GetItemInput{
		TableName: awssdk.String(name),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Item, nil
}

// DynamoPutItem put an item
func DynamoPutItem(name string, items map[string]*dynamodb.AttributeValue) (result map[string]*dynamodb.AttributeValue, e error) {
	resp, err := dynamodb.New(session.New(), dynamoCfg).PutItem(&dynamodb.PutItemInput{
		TableName: awssdk.String(name),
		Item:      items,
	})
	if err != nil {
		return nil, err
	}
	return resp.Attributes, nil
}

// DynamoDeleteItem put an item
func DynamoDeleteItem(name string, attributes map[string]string) (result map[string]*dynamodb.AttributeValue, e error) {
	items := map[string]*dynamodb.AttributeValue{}
	for key, value := range attributes {
		items[key] = &dynamodb.AttributeValue{
			S: awssdk.String(value),
		}
	}
	resp, err := dynamodb.New(session.New(), dynamoCfg).DeleteItem(&dynamodb.DeleteItemInput{
		TableName: awssdk.String(name),
		Key:       items,
	})
	if err != nil {
		return nil, err
	}
	return resp.Attributes, nil
}

// DynamoCreateTable creates a dynamodb table
func DynamoCreateTable(name string, attributes map[string]string, keys map[string]string,
	readCapacityUnits int64, writeCapacityUnits int64) (table *dynamodb.TableDescription, e error) {
	attributeDefinitions := []*dynamodb.AttributeDefinition{}
	for key, value := range attributes {
		attributeDefinitions = append(attributeDefinitions, &dynamodb.AttributeDefinition{
			AttributeName: awssdk.String(key),
			AttributeType: awssdk.String(value),
		})
	}
	keySchema := []*dynamodb.KeySchemaElement{}
	for key, value := range keys {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: awssdk.String(key),
			KeyType:       awssdk.String(value),
		})
	}
	throughput := &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  awssdk.Int64(readCapacityUnits),
		WriteCapacityUnits: awssdk.Int64(writeCapacityUnits),
	}
	resp, err := dynamodb.New(session.New(), dynamoCfg).CreateTable(&dynamodb.CreateTableInput{
		TableName:             awssdk.String(name),
		AttributeDefinitions:  attributeDefinitions,
		KeySchema:             keySchema,
		ProvisionedThroughput: throughput,
	})
	if err != nil {
		return nil, err
	}
	return resp.TableDescription, nil
}

// DynamoDropTable drops a specific dynamodb table
func DynamoDropTable(name string) (table *dynamodb.TableDescription, e error) {
	resp, err := dynamodb.New(session.New(), dynamoCfg).DeleteTable(&dynamodb.DeleteTableInput{
		TableName: awssdk.String(name),
	})
	if err != nil {
		return nil, err
	}
	return resp.TableDescription, nil
}
