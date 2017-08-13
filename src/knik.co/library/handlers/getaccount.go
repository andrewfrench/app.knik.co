package handlers

import (
	"log"
	"knik.co/library/account/instagram"
)

type GetAccountRequest struct {
	Token string
	AccountId string
}

func GetAccountHandler(req *GetAccountRequest) map[string]interface{} {
	log.Printf("Entering GetAccountHandler")
	defer log.Printf("Exiting GetAccountHandler")

	_, u, resp := EnsureAuthentication(req.Token)
	if len(resp) > 0 {
		return resp
	}

	a, err := instagram.GetAccountById(req.AccountId)
	if err != nil {
		log.Printf("Error getting account: %s", err.Error())
		return map[string]interface{}{
			"error": "Account not found",
			"bounce": true,
		}
	}

	ownedByUser := a.OwnerId == u.Id
	log.Printf("Account owned by user: %t", ownedByUser)

	if !(u.Admin || ownedByUser) {
		log.Printf("User does not have credentials to view this account")
		return map[string]interface{}{
			"error": "You cannot view this account",
			"bounce": true,
		}
	}

	a.RefreshIfStale()

	return map[string]interface{}{
		"account": a,
		"owned_by_user": ownedByUser,
	}
}
