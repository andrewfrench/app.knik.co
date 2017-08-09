package handlers

import (
	"log"
	"knik.co/library/session"
	"knik.co/library/user"
	"knik.co/library/account/instagram"
)

type GetAccountRequest struct {
	Token string
	AccountId string
}

func GetAccountHandler(req *GetAccountRequest) map[string]interface{} {
	log.Printf("Entering GetAccountHandler")
	defer log.Printf("Exiting GetAccountHandler")

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

	ownedByUser := a.OwnerId == u.Id
	log.Printf("Account owned by user: %t", ownedByUser)

	a.RefreshIfStale()

	return map[string]interface{}{
		"account": a,
		"owned_by_user": ownedByUser,
	}
}
