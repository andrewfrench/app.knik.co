package handlers

import (
	"log"
)

type SignOutRequest struct {
	Token string
}

func SignOutHandler(req *SignOutRequest) map[string]interface{} {
	log.Printf("Entering SignOutHandler")
	defer log.Printf("Exiting SignOutHandler")

	s, _, resp := EnsureAuthentication(req.Token)
	if len(resp) > 0 {
		return resp
	}

	err := s.Delete()
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
