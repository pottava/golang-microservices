package aws

/**
 * @see https://github.com/aws/aws-sdk-go/blob/master/service/dynamodb/api.go
 */

import (
	"strconv"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoS gets string from AttributeValue
func DynamoS(data map[string]*dynamodb.AttributeValue, key string) string {
	if value, ok := data[key]; ok {
		return *value.S
	}
	return "" 
}

// DynamoN gets int from AttributeValue
func DynamoN(data map[string]*dynamodb.AttributeValue, key string) int {
	if value, ok := data[key]; ok {
		i, _ := strconv.Atoi(*value.N)
		return i
	}
	return 0
}

// DynamoN64 gets int64 from AttributeValue
func DynamoN64(data map[string]*dynamodb.AttributeValue, key string) int64 {
	if value, ok := data[key]; ok {
		i64, _ := strconv.ParseInt(*value.N, 10, 64)
		return i64
	}
	return 0
}

// DynamoD gets time.Time from AttributeValue
func DynamoD(data map[string]*dynamodb.AttributeValue, key string) time.Time {
	if value, ok := data[key]; ok {
		i64, _ := strconv.ParseInt(*value.N, 10, 64)
		return time.Unix(i64, 0)
	}
	return time.Now()
}

// DynamoB gets boolean from AttributeValue
func DynamoB(data map[string]*dynamodb.AttributeValue, key string) bool {
	if value, ok := data[key]; ok {
		return *value.BOOL
	}
	return false
}

// DynamoAttributeS makes string to DynamoDB String Attribute
func DynamoAttributeS(value string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		S: awssdk.String(value),
	}
}

// DynamoAttributeN makes int to DynamoDB Number Attribute
func DynamoAttributeN(value int) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		N: awssdk.String(strconv.Itoa(value)),
	}
}

// DynamoAttributeN64 makes int64 to DynamoDB Number Attribute
func DynamoAttributeN64(value int64) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		N: awssdk.String(strconv.FormatInt(value, 10)),
	}
}

// DynamoAttributeD makes time.Time to DynamoDB Number Attribute
func DynamoAttributeD(value time.Time) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		N: awssdk.String(strconv.FormatInt(value.Unix(), 10)),
	}
}

// DynamoAttributeB makes boolean to DynamoDB Number Attribute
func DynamoAttributeB(value bool) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		BOOL: awssdk.Bool(value),
	}
}
