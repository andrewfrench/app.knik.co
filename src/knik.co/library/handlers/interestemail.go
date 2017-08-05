package handlers

import (
	"log"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"time"
	"knik.co/library/database"
)

type InterestEmailRequest struct {
	Email string
}

const table string = "knik.co-interest-emails"

func InterestEmailHander(req *InterestEmailRequest) map[string]interface{} {
	log.Printf("Entering InterestEmailHandler")
	defer log.Printf("Exiting InterestEmailHandler")

	params := &dynamodb.PutItemInput{
		TableName: &table,
		Item: map[string]*dynamodb.AttributeValue{
			"email": {S: &req.Email},
			"timestamp": {N: aws.String(fmt.Sprintf("%s", time.Now().Unix()))},
		},
	}

	_, err := database.PutItem(params)
	if err != nil {
		log.Printf("Error putting email: %s", err.Error())
		return map[string]interface{}{
			"error": "Error saving email",
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}
