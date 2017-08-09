package handlers

import (
	"log"
	"knik.co/library/session"
	"knik.co/library/user"
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

	s, err := session.GetSessionBySessionId(req.Token)
	if err != nil {
		log.Printf("Error getting session: %s", err.Error())
		return map[string]interface{}{
			"error": "Unauthenticated",
		}
	}

	u, err := user.GetUserById(s.UserId)
	if err != nil {
		log.Printf("Error getting user: %s", err.Error())
		return map[string]interface{}{
			"error": "Unauthenticated",
		}
	}

	a, err := instagram.GetAccountById(req.AccountId)
	if err != nil {
		log.Printf("Error getting account: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to update account",
		}
	}

	if u.Id != a.OwnerId {
		log.Printf("User does not own account")
		return map[string]interface{}{
			"error": "Unauthenticated",
		}
	}

	newAccount, err := account.Get(req.Username)
	if err != nil {
		log.Printf("Error getting account from new username: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to update account",
		}
	}

	if a.InstagramId != newAccount.Id {
		log.Printf("Instagram ID mismatch: %s != %s", a.InstagramId, newAccount.Id)
		return map[string]interface{}{
			"error": "Unable to update account",
		}
	}

	a.Username = newAccount.Username
	err = a.Insert()
	if err != nil {
		log.Printf("Error updating account: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to update account",
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}
