package handlers

import (
	"log"
	"knik.co/library/session"
	"knik.co/library/user"
)

type GetSelfRequest struct {
	Token string
}

func GetSelfHandler(req *GetSelfRequest) map[string]interface{} {
	log.Printf("Entering GetSelfHandler")
	defer log.Printf("Exiting GetSelfHandler")

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

	return map[string]interface{}{
		"profile": u,
	}
}
