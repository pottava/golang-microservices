package aws

import (
	"reflect"
	"strconv"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestDynamoS(t *testing.T) {
	actual := DynamoS(map[string]*dynamodb.AttributeValue{
		"ID": &dynamodb.AttributeValue{
			S: awssdk.String("123"),
		},
	}, "ID")
	expected := "123"
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = DynamoS(map[string]*dynamodb.AttributeValue{
		"ID": &dynamodb.AttributeValue{
			S: awssdk.String("123"),
		},
	}, "foo")
	expected = ""
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestDynamoN(t *testing.T) {
	actual := DynamoN(map[string]*dynamodb.AttributeValue{
		"TypeId": &dynamodb.AttributeValue{
			N: awssdk.String(strconv.Itoa(123)),
		},
	}, "TypeId")
	expected := 123
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
	actual = DynamoN(map[string]*dynamodb.AttributeValue{
		"TypeId": &dynamodb.AttributeValue{
			N: awssdk.String(strconv.Itoa(123)),
		},
	}, "foo")
	expected = 0
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestDynamoAttributeS(t *testing.T) {
	actual := DynamoAttributeS("123")
	expected := &dynamodb.AttributeValue{
		S: awssdk.String("123"),
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestDynamoAttributeN(t *testing.T) {
	actual := DynamoAttributeN(123)
	expected := &dynamodb.AttributeValue{
		N: awssdk.String(strconv.Itoa(123)),
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}
