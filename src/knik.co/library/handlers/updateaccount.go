package handlers

import (
	"log"
	"knik.co/library/account/instagram"
	"knik.co/library/session"
	"knik.co/library/user"
)

type UpdateAccountRequest struct {
	Token string
	AccountId string
	Market string
	Location string
	Experience string
	Summary string
}

func UpdateAccountHandler(req *UpdateAccountRequest) map[string]interface{} {
	log.Printf("Entering UpdateAccountHandler")
	defer log.Printf("Exiting UpdateAccountHandler")

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
			"error": "Account not found",
		}
	}

	if a.OwnerId != u.Id {
		log.Printf("User does not own account")
		return map[string]interface{}{
			"error": "You do not own this account",
		}
	}

	a.Market = req.Market
	a.Location = req.Location
	a.Experience = req.Experience
	a.Summary = req.Summary
	if a.Insert() != nil {
		log.Printf("Error updating account: %s", err.Error())
		return map[string]interface{}{
			"error": "Error updating account",
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}