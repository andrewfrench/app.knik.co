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
	if UserHasValidSession(userId) {
		log.Println("User has a valid session")
		sess, _ := GetSessionByUserId(userId)
		sess.Delete()
	}

	sessionId := random.RandomString(32)
	for database.Exists("session_id", sessionId, table()) {
		sessionId = random.RandomString(32)
	}

	return &Session{
		SessionId: sessionId,
		UserId: userId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

func UserHasValidSession(userId string) bool {
	log.Println("Checking if user has a valid session")

	sess, err := GetSessionByUserId(userId)
	if err != nil {
		return false
	}

	if sess.ExpiresAt.Before(time.Now()) {
		log.Println("Session exists, but has expired")
		sess.Delete()
		return false
	}

	return true
}

func Authenticate(sessionId, userId string) bool {
	log.Printf("Authenticating session with id: %s and user with id: %s", sessionId, userId)

	sess, err := GetSessionBySessionId(sessionId)
	if err != nil {
		return false
	}

	if sess.ExpiresAt.Before(time.Now()) {
		return false
	}

	if sess.UserId != userId {
		return false
	}

	return true
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

	s := responseItemToSession(resp.Item)
	if s.SessionId == "" {
		return &Session{}, errors.New("No session exists")
	}

	return s, err
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

func GetSessionByUserId(userId string) (*Session, error) {
	log.Printf("Checking for existing session with userId %s", userId)

	params := &dynamodb.QueryInput{
		IndexName: aws.String("user_id-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user_id": {S: aws.String(userId)},
		},
		KeyConditionExpression: aws.String("user_id = :user_id"),
		TableName: table(),
	}

	resp, err := database.Query(params)
	if err != nil {
		log.Printf("Error getting session by user id: %s", err.Error())
		return &Session{}, err
	}

	if *resp.Count == 0 {
		return &Session{}, err
	}

	return responseItemToSession(resp.Items[0]), err
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
