package user

import (
	"log"
	"time"
	"knik.co/library/database"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/andrewfrench/random"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"errors"
	"knik.co/library/account/instagram"
	"knik.co/library/common"
)

type User struct {
	// General profile information
	Id    string `json:"user_id"`
	Email string `json:"user_email"`
	CreatedAt int64 `json:"created_at"`
	password string

	// Accounts
	Accounts []instagram.Account `json:"accounts"`
}

const table string = "knik.co-users"

type igApiResponse struct {
	Data struct {
		Id             string `json:"id"`
		Username       string `json:"username"`
		FullName       string `json:"full_name"`
		ProfilePicture string `json:"profile_picture"`
		Bio            string `json:"bio"`
		Website        string `json:"website"`
		Counts struct {
			Media      int `json:"media"`
			Follows    int `json:"follows"`
			FollowedBy int `json:"followed_by"`
		} `json:"counts"`
	} `json:"data"`
}

func Create(email, password string) *User {
	log.Println("Creating user struct")

	log.Println("Generating user ID")
	id := random.RandomString(10)
	for database.Exists("user_id", id, table) {
		log.Printf("User ID %s already allocated, generating another...", id)
		id = random.RandomString(10)
	}

	return &User{
		Id: id,
		Email: email,
		password: common.Hash(password),
		CreatedAt: time.Now().Unix(),
	}
}

func GetUserByEmail(email string) (*User, error) {
	log.Printf("Getting user with email %s", email)

	params := &dynamodb.QueryInput{
		IndexName: aws.String("user_email-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {S: aws.String(email)},
		},
		KeyConditionExpression: aws.String("user_email=:email"),
		TableName: aws.String(table),
	}

	resp, err := database.Query(params)
	if err != nil {
		return &User{}, err
	}

	if *resp.Count == 0 {
		return &User{}, err
	}

	u, err := responseItemToUser(resp.Items[0])
	return u, err
}

func EmailIsRegistered(email string) bool {
	log.Println("Checking if email is registered")

	u, _ := GetUserByEmail(email)

	return u.Id != ""
}

func Authenticate(email, password string) (*User, error) {
	log.Println("Authenticating credentials")

	u, err := GetUserByEmail(email)
	if err != nil {
		log.Printf("Error getting user: %s", err.Error())
		return &User{}, err
	}

	if u.password != common.Hash(password) {
		log.Println("Bad password")
		return &User{}, errors.New("Bad password")
	}

	return u, err
}

func (u *User) Insert() error {
	log.Println("Inserting user")

	params := &dynamodb.PutItemInput{
		Item: u.AttributeValues(),
		TableName: aws.String(table),
	}

	_, err := database.PutItem(params)

	return err
}

func GetUserById(id string) (*User, error) {
	log.Printf("Getting user with id: %s", id)

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(table),
	}

	resp, err := database.GetItem(params)
	if err != nil {
		return &User{}, err
	}

	u, err := responseItemToUser(resp.Item)
	u.Accounts, err = instagram.GetAccountsByOwnerId(u.Id)
	return u, err
}

func (u *User) UpdateExistingUser() error {
	params := &dynamodb.PutItemInput{
		Item: u.AttributeValues(),
		TableName: aws.String(table),
	}

	_, err := database.PutItem(params)

	return err
}
