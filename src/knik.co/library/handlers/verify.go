package handlers

import (
	"log"
	"knik.co/library/user"
	"knik.co/library/session"
	"knik.co/library/account/instagram"
)

type VerifyRequest struct {
	Token string
	Username string
	AuthCode string
}

func VerifyHandler(req *VerifyRequest) map[string]interface{} {
	log.Printf("Entering VerifyHandler")
	defer log.Printf("Exiting VerifyHandler")

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

	// TODO: Get recent media cache here, move cache validity logic to account getter
	// TODO: Bust out checks to ensure account isn't verified to their own methods, do that check here
	acc := instagram.Create(u.Id, req.Username)
	err = acc.Verify()
	if err != nil {
		log.Printf("Error verifying account: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to verify account",
		}
	}

	return map[string]interface{}{
		"success": true,
		"account_id": acc.AccountId,
	}
}
