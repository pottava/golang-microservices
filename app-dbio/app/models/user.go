package models

import (
	"sort"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/pottava/golang-microservices/app-dbio/app/aws"
	"github.com/pottava/golang-microservices/app-dbio/app/logs"
)

const userTable = "gomicroservices-users"

// User represents user's user
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Users is a type of User slice
type Users []*User

var userOnce sync.Once

func init() {
	userOnce.Do(func() {
		comfirmUserTableExists()
	})
}

// GetUsers lists all users from DynamoDB
//  @return users []models.User
func GetUsers() (users Users, count int64, err error) {
	records, count, err := aws.DynamoScan(userTable)
	if err != nil {
		return users, 0, err
	}
	return toUsers(records), count, nil
}

// GetUser retrives a specified user from DynamoDB
//  @param  id string
//  @return user models.User
func GetUser(id string) (user *User, found bool) {
	record, err := aws.DynamoRecord(userTable, map[string]*dynamodb.AttributeValue{
		"ID": aws.DynamoAttributeS(id),
	})
	if (err != nil) || len(record) == 0 {
		return nil, false
	}
	return toUser(record), true
}

// cast DynamoDB records to Users
func toUsers(records []map[string]*dynamodb.AttributeValue) (users Users) {
	for _, record := range records {
		user := toUser(record)
		users = append(users, user)
	}
	if len(users) == 0 {
		users = make([]*User, 0)
	}
	sort.Sort(users)
	return users
}

// cast DynamoDB record to a User
func toUser(record map[string]*dynamodb.AttributeValue) *User {
	user := User{}
	user.ID = aws.DynamoS(record, "ID")
	user.Name = aws.DynamoS(record, "Name")
	return &user
}

// Persist persists its state
func (u *User) Persist() error {
	items := map[string]*dynamodb.AttributeValue{}
	items["ID"] = aws.DynamoAttributeS(u.ID)
	items["Name"] = aws.DynamoAttributeS(u.Name)
	_, err := aws.DynamoPutItem(userTable, items)
	if err != nil {
		logs.Error.Printf("User#Persist. Items: %v, Error: %v", items, err)
	}
	return err
}

func comfirmUserTableExists() (found bool, err error) {
	if _, err := aws.DynamoTable(userTable); err == nil {
		return true, nil
	}
	logs.Debug.Print("[model] User table was not found. Try to make it. @aws.DynamoCreateTable")
	attributes := map[string]string{
		"ID": "S",
	}
	keys := map[string]string{
		"ID": "HASH",
	}
	_, err = aws.DynamoCreateTable(userTable, attributes, keys, 1, 1)
	if err != nil {
		logs.Error.Printf("DynamoCreateTable. Name: %v, Error: %v", userTable, err)
	}
	return false, err
}

func (s Users) Len() int {
	return len(s)
}

func (s Users) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Users) Less(i, j int) bool {
	a, _ := strconv.Atoi(s[i].ID)
	b, _ := strconv.Atoi(s[j].ID)
	return a < b
}
