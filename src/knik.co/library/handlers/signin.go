package handlers

import (
	"log"
	"knik.co/library/user"
	"knik.co/library/session"
)

type SignInRequest struct {
	Email string
	Password string
}

func SignInHandler(req *SignInRequest) map[string]interface{} {
	log.Printf("Entering SignInHandler")
	defer log.Printf("Exiting SignInHandler")

	u, err := user.Authenticate(req.Email, req.Password)
	if err != nil {
		return map[string]interface{}{
			"error": "Invalid credentials",
		}
	}

	s := session.Create(u.Id)
	err = s.Insert()
	if err != nil {
		return map[string]interface{}{
			"error": "Error during authentication",
		}
	}

	return map[string]interface{}{
		"success": true,
		"token": s.SessionId,
	}
}
