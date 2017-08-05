package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

var db *dynamodb.DynamoDB

func init() {
	log.Println("Initializing database session")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		panic(err.Error())
	}

	db = dynamodb.New(sess)
}

func GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	log.Println("Getting item")
	return db.GetItem(input)
}

func PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	log.Println("Putting item")
	return db.PutItem(input)
}

func UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	log.Println("Updating item")
	return db.UpdateItem(input)
}

func DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	log.Println("Deleting item")
	return db.DeleteItem(input)
}

func Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	log.Println("Scanning items")
	return db.Scan(input)
}

func Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	log.Println("Querying items")
	return db.Query(input)
}

func Exists(idField, id string, table *string) bool {
	log.Println("Checking if item exists")
	params := &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			idField: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{S: aws.String(id)},
				},
			},
		},
		TableName: table,
	}

	resp, err := Query(params)
	if err != nil {
		panic(err.Error())
	}

	return len(resp.Items) > 0
}
