package handlers

import (
	"log"
	"knik.co/library/user"
	"fmt"
	"knik.co/library/session"
)

type SignUpRequest struct {
	Email string
	Password string
}

func SignUpHandler(req *SignUpRequest) map[string]interface{} {
	log.Printf("Entering SignUpHandler")
	defer log.Printf("Exiting SignUpHandler")

	minPasswordLength := 6
	if len(req.Password) < minPasswordLength {
		log.Printf("Password too short")
		return map[string]interface{}{
			"error": fmt.Sprintf("Password must be %d characters long", minPasswordLength),
		}
	}

	if user.EmailIsRegistered(req.Email) {
		log.Printf("Email %s already registered", req.Email)
		return map[string]interface{}{
			"error": "Email is already registered",
		}
	}

	u := user.Create(req.Email, req.Password)
	err := u.Insert()
	if err != nil {
		log.Printf("Error inserting user: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to create user",
		}
	}

	s := session.Create(u.Id)
	err = s.Insert()
	if err != nil {
		log.Printf("Error inserting session: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to create user",
		}
	}

	return map[string]interface{}{
		"success": true,
		"token": s.SessionId,
	}
}