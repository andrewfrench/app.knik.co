package session

import (
	"time"
	"github.com/andrewfrench/random"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"knik.co/library/database"
	"fmt"
	"log"
	"strconv"
	"errors"
	"os"
)

type Session struct {
	SessionId string
	UserId string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func table() *string {
	s := os.Getenv("TABLE_SESSIONS")
	return &s
}

func Create(userId string) *Session {
	id := random.RandomString(32)
	for database.Exists("session_id", id, table()) {
		id = random.RandomString(32)
	}

	return &Session{
		SessionId: id,
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

func GetSessionBySessionId(sessionId string) (*Session, error) {
	log.Printf("Getting session with ID: %s", sessionId)

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"session_id": {S: aws.String(sessionId)},
		},
		TableName: table(),
	}

	resp, err := database.GetItem(params)
	if err != nil {
		log.Printf("Error getting session by session ID: %s", err.Error())
		return &Session{}, err
	}

	if len(resp.Item) == 0 {
		log.Printf("No such session exists")
		return &Session{}, errors.New("No session exists")
	}

	return responseItemToSession(resp.Item), err
}

func (s *Session) Delete() error {
	log.Println("Deleting session")

	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"session_id": {S: aws.String(s.SessionId)},
		},
		TableName: table(),
	}

	_, err := database.DeleteItem(params)
	if err != nil {
		log.Printf("Error while deleting session: %s", err.Error())
	}

	return err
}

func (s *Session) Insert() error {
	log.Println("Inserting session")

	params := &dynamodb.PutItemInput{
		Item: s.AttributeValues(),
		TableName: table(),
	}

	_, err := database.PutItem(params)

	return err
}

func (s *Session) AttributeValues() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"session_id": {S: aws.String(s.SessionId)},
		"user_id": {S: aws.String(s.UserId)},
		"session_created_at": {N: aws.String(fmt.Sprintf("%d", s.CreatedAt.Unix()))},
		"session_expires_at": {N: aws.String(fmt.Sprintf("%d", s.ExpiresAt.Unix()))},
	}
}

func responseItemToSession(item map[string]*dynamodb.AttributeValue) *Session {
	log.Println("Building session from response item")

	log.Println("Casting strings to ints")
	createdAtInt, err := strconv.Atoi(*item["session_created_at"].N)
	if err != nil {
		log.Fatalf("Failed to convert to int: %s", err.Error())
	}

	expiresAtInt, err := strconv.Atoi(*item["session_expires_at"].N)
	if err != nil {
		log.Fatalf("Failed to convert to int: %s", err.Error())
	}

	log.Println("Building and returning session struct")
	return &Session{
		SessionId: *item["session_id"].S,
		UserId: *item["user_id"].S,
		CreatedAt: time.Unix(int64(createdAtInt), 0),
		ExpiresAt: time.Unix(int64(expiresAtInt), 0),
	}
}
