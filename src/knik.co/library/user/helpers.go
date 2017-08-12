package user

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
	"strconv"
)

func (u *User) AttributeValues() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		// General user data
		"user_id": {S: aws.String(u.Id)},
		"user_email": {S: aws.String(u.Email)},
		"user_created_at": {S: aws.String(fmt.Sprintf("%d", u.CreatedAt))},
		"user_password": {S: aws.String(u.Password)},
		"is_admin": {BOOL: &u.Admin},
	}
}

func responseItemToUser(item map[string]*dynamodb.AttributeValue) (*User, error) {
	createdAt, err := strconv.Atoi(*item["user_created_at"].S)
	if err != nil {
		return &User{}, err
	}

	return &User{
		Id:        *item["user_id"].S,
		Email:     *item["user_email"].S,
		CreatedAt: int64(createdAt),
		Admin:     *item["is_admin"].BOOL,
		Password:  *item["user_password"].S,
	}, err
}
