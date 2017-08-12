package handlers

import (
	"log"
	"knik.co/library/account/instagram"
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

	_, u, resp := EnsureAuthentication(req.Token)
	if len(resp) > 0 {
		return resp
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