package instagram

import (
	"github.com/andrewfrench/random"
	"knik.co/library/database"
	"log"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"strconv"
	"errors"
	"time"
	ig "github.com/andrewfrench/instagram-api-bypass/account"
	"knik.co/library/account"
	"encoding/json"
)

type Account struct {
	AccountId        string `json:"id"`
	InstagramId      string
	OwnerId          string
	IsVerified       bool   `json:"verified"`
	VerifiedAt       int64
	CreatedAt        int64  `json:"created_at"`
	UpdatedAt        int64  `json:"updated_at"`
	VerificationCode string
	Username         string `json:"username"`
	Followers        int    `json:"followers"`
	Url              string `json:"url"`
	ProfilePicUrl	 string `json:"profile_pic_url"`
	RecentImageUrl   string `json:"recent_image_url"`
	RecentMedia      []ig.RecentMedia `json:"recent_media"`
	Market           string `json:"market"`
	Location         string `json:"location"`
	Experience       string `json:"experience"`
	Summary          string `json:"summary"`
	AverageInteractions float32 `json:"average_interactions"`
	PostPeriod		 int64 `json:"post_period"`
}

const table string = "knik.co-accounts-instagram"

func Create(ownerId, username string) *Account {
	log.Println("Creating account")

	randomId := random.RandomString(10)
	for database.Exists("account_id", randomId, table) {
		randomId = random.RandomString(10)
	}

	return &Account{
		AccountId:   randomId,
		OwnerId:     ownerId,
		Username:    username,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		RecentMedia: []ig.RecentMedia{},
	}
}

func GetAccountById(accountId string) (*Account, error) {
	log.Printf("Getting account with ID: %s", accountId)

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"account_id": {S:aws.String(accountId)},
		},
		TableName: aws.String(table),
	}

	resp, err := database.GetItem(params)
	if err != nil {
		return &Account{}, err
	}

	acc := responseItemToAccount(resp.Item)
	acc.RefreshIfStale()
	return &acc, err
}

func GetAccountsByOwnerId(ownerId string) ([]Account, error) {
	log.Printf("Getting accounts with owner: %s", ownerId)

	params := &dynamodb.QueryInput{
		IndexName: aws.String("owner_id-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":owner_id": {S: aws.String(ownerId)},
		},
		KeyConditionExpression: aws.String("owner_id = :owner_id"),
		TableName: aws.String(table),
	}

	resp, err := database.Query(params)
	if err != nil {
		log.Printf("Failed to query accounts: %s", err.Error())
	}

	accs := []Account{}

	for _, i := range resp.Items {
		accs = append(accs, responseItemToAccount(i))
	}

	return accs, err
}

func AccountIdExists(instagramId string) bool {
	log.Printf("Querying instagram accounts for id: %s", instagramId)

	params := &dynamodb.QueryInput{
		IndexName: aws.String("instagram_id-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":instagram_id": {S: aws.String(instagramId)},
		},
		KeyConditionExpression: aws.String("instagram_id = :instagram_id"),
		TableName: aws.String(table),
	}

	resp, err := database.Query(params)
	if err != nil {
		log.Fatalf("Failed to query existing accounts")
	}

	return *resp.Count > 0
}

func (a *Account) Insert() error {
	log.Println("Inserting account")

	params := &dynamodb.PutItemInput{
		Item: a.attributeValues(),
		TableName: aws.String(table),
	}

	_, err := database.PutItem(params)
	if err != nil{
		log.Printf("Failed to insert account: %s", err.Error())
	}

	return err
}

func (a *Account) RefreshIfStale() {
	log.Printf("Entering account.RefreshIfStale()")
	defer log.Printf("Exiting account.RefreshIfStale()")

	if time.Now().After(time.Unix(a.UpdatedAt, 0).Add(time.Hour)) {
		log.Printf("Account information is stale")
		a.Refresh()
	} else {
		log.Printf("Account information is not stale")
	}
}

func (a *Account) Refresh() {
	log.Printf("Entering account.Refresh()")
	defer log.Printf("Exiting account.Refresh()")

	acc, err := ig.Get(a.Username)
	if err != nil {
		log.Printf("Error getting account: %s", err.Error())
	}

	a.Followers = acc.Followers
	a.UpdatedAt = time.Now().Unix()
	a.RecentMedia = acc.RecentMedia
	a.ProfilePicUrl = acc.ProfilePicUrl
	a.Insert()
}

func (a *Account) Verify() error {
	log.Println("Verifying Instagram account")

	acc, err := ig.Get(a.Username)
	if err != nil {
		return err
	}

	if AccountIdExists(acc.Id) {
		return errors.New("Account already verified")
	}

	authCodeCandidates := account.ExtractCodeCandidates(acc.Biography)
	matchFound := false
	actualCode := account.CodeGen(a.Username, a.OwnerId)
	log.Printf("Looking for code: %s", actualCode)

	for _, candidate := range authCodeCandidates {
		log.Printf("Inspecting %s", candidate)
		if candidate == actualCode {
			matchFound = true
			break
		}
	}

	if !matchFound {
		return errors.New("Verification code not found")
	}

	a.IsVerified = true
	a.InstagramId = acc.Id
	a.VerifiedAt = time.Now().Unix()
	a.UpdatedAt = time.Now().Unix()
	a.Url = fmt.Sprintf("https://www.instagram.com/%s/", acc.Username)
	a.Followers = acc.Followers
	a.VerificationCode = actualCode
	a.RecentImageUrl = acc.RecentMedia[0].ThumbnalSrc
	a.RecentMedia = acc.RecentMedia
	a.ProfilePicUrl = acc.ProfilePicUrl
	return a.Insert()
}

func (a *Account) attributeValues() map[string]*dynamodb.AttributeValue {
	log.Println("Building account attribute value map")

	avm := map[string]*dynamodb.AttributeValue{
		"account_id": {S: aws.String(a.AccountId)},
		"owner_id": {S: aws.String(a.OwnerId)},
		"username": {S: aws.String(a.Username)},
		"url": {S: aws.String(a.Url)},
		"is_verified": {BOOL: aws.Bool(a.IsVerified)},
		"followers": {N: aws.String(fmt.Sprintf("%d", a.Followers))},
		"verification_code": {S: aws.String(a.VerificationCode)},
		"verified_at": {N: aws.String(fmt.Sprintf("%d", a.VerifiedAt))},
		"created_at": {N: aws.String(fmt.Sprintf("%d", a.CreatedAt))},
		"updated_at": {N: aws.String(fmt.Sprintf("%d", a.UpdatedAt))},
		"recent_image_url": {S: aws.String(a.RecentImageUrl)},
		"instagram_id": {S: aws.String(a.InstagramId)},
	}

	if len(a.ProfilePicUrl) > 0 {
		avm["profile_pic_url"] = &dynamodb.AttributeValue{S: aws.String(a.ProfilePicUrl)}
	}

	if len(a.Market) > 0 {
		avm["market"] = &dynamodb.AttributeValue{S: aws.String(a.Market)}
	}

	if len(a.Location) > 0 {
		avm["location"] = &dynamodb.AttributeValue{S: aws.String(a.Location)}
	}

	if len(a.Experience) > 0 {
		avm["experience"] = &dynamodb.AttributeValue{S: aws.String(a.Experience)}
	}

	if len(a.Summary) > 0 {
		avm["summary"] = &dynamodb.AttributeValue{S: aws.String(a.Summary)}
	}

	if a.RecentMedia != nil && len(a.RecentMedia) > 0 {
		marshalled, _ := json.Marshal(a.RecentMedia)
		avm["recent_media"] = &dynamodb.AttributeValue{B: marshalled}
	}

	return avm
}

func responseItemToAccount(item map[string]*dynamodb.AttributeValue) Account {
	log.Println("Building account from response item")

	log.Println("Converting strings to integers")
	followers, _ := strconv.Atoi(*item["followers"].N)
	verifiedAtInt, _ := strconv.Atoi(*item["verified_at"].N)
	createdAtInt, _ := strconv.Atoi(*item["created_at"].N)
	updatedAtInt, _ := strconv.Atoi(*item["updated_at"].N)

	log.Println("Building account struct")
	a := Account{
		AccountId:        *item["account_id"].S,
		OwnerId:          *item["owner_id"].S,
		Username:         *item["username"].S,
		Url:              *item["url"].S,
		IsVerified:       *item["is_verified"].BOOL,
		Followers:        followers,
		VerifiedAt:       int64(verifiedAtInt),
		VerificationCode: *item["verification_code"].S,
		CreatedAt:        int64(createdAtInt),
		UpdatedAt:        int64(updatedAtInt),
		RecentImageUrl:   *item["recent_image_url"].S,
		InstagramId:	  *item["instagram_id"].S,
	}

	if f, e := item["profile_pic_url"]; e{
		a.ProfilePicUrl = *f.S
	}

	if f, e := item["market"]; e {
		a.Market = *f.S
	}

	if f, e := item["location"]; e {
		a.Location = *f.S
	}

	if f, e := item["experience"]; e {
		a.Experience = *f.S
	}

	if f, e := item["summary"]; e {
		a.Summary = *f.S
	}

	if f, e := item["recent_media"]; e {
		err := json.Unmarshal(f.B, &a.RecentMedia)
		if err != nil {
			log.Printf("Error unmarshalling cached recent media: %s", err.Error())
		}
	}

	if len(a.RecentMedia) > 0 {
		interactions := 0
		for _, m := range a.RecentMedia {
			interactions += m.Likes
			interactions += m.Comments
		}

		a.AverageInteractions = float32(interactions) / float32(len(a.RecentMedia))

		timeDiff := a.RecentMedia[0].Date - a.RecentMedia[len(a.RecentMedia) - 1].Date
		a.PostPeriod = int64(timeDiff) / int64(len(a.RecentMedia))
	}

	return a
}
