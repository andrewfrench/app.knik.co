package handlers

import (
	"log"
	"knik.co/library/account/instagram"
	"github.com/andrewfrench/instagram-api-bypass/account"
	"strings"
)

type UpdateUsernameRequest struct {
	Token string
	AccountId string
	Username string
}

func UpdateUsernameHandler(req *UpdateUsernameRequest) map[string]interface{} {
	log.Printf("Entering UpdateUsernameHandler")
	defer log.Printf("Exiting UpdateUsernameHandler")

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

	_, u, resp := EnsureAuthentication(req.Token)
	if len(resp) > 0 {
		return resp
	}

	a, err := instagram.GetAccountById(req.AccountId)
	if err != nil {
		log.Printf("Error getting account: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to update username",
		}
	}

	if u.Id != a.OwnerId {
		log.Printf("User does not own account")
		return map[string]interface{}{
			"error": "Unauthenticated",
			"bounce": true,
		}
	}

	if a.Username == req.Username {
		log.Printf("New and current username are equal")
		return map[string]interface{}{
			"error": "This is already your username",
		}
	}

	newAccount, err := account.Get(req.Username)
	if err != nil {
		log.Printf("Error getting account from new username: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to update username",
		}
	}

	if a.InstagramId != newAccount.Id {
		log.Printf("Instagram ID mismatch: %s != %s", a.InstagramId, newAccount.Id)
		return map[string]interface{}{
			"error": "This is a different account",
		}
	}

	a.Username = newAccount.Username
	err = a.Insert()
	if err != nil {
		log.Printf("Error updating account: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to update username",
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}
