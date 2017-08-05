package handlers

import (
	"log"
	"knik.co/library/session"
	"knik.co/library/user"
	"knik.co/library/account"
)

type GetAuthCodeRequest struct {
	Token string
	Username string
}

func GetAuthCodeHandler(req *GetAuthCodeRequest) map[string]interface{} {
	log.Printf("Entering GetAuthCodeHandler")
	defer log.Printf("Exiting GetAuthCodeHandler")

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

	code := account.CodeGen(req.Username, u.Id)

	return map[string]interface{}{
		"auth_code": code,
	}
}
