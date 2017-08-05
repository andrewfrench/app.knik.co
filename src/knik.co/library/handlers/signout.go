package handlers

import (
	"log"
	"knik.co/library/session"
)

type SignOutRequest struct {
	Token string
}

func SignOutHandler(req *SignOutRequest) map[string]interface{} {
	log.Printf("Entering SignOutHandler")
	defer log.Printf("Exiting SignOutHandler")

	s, err := session.GetSessionBySessionId(req.Token)
	if err != nil {
		log.Printf("Error getting session: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to sign out",
		}
	}

	err = s.Delete()
	if err != nil {
		log.Printf("Error deleting session: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to sign out",
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}
