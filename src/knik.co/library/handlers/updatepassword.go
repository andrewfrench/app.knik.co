package handlers

import (
	"log"
	"knik.co/library/user"
	"fmt"
)

type UpdatePasswordRequest struct {
	Token string
	Password string
	NewPassword string
}

func UpdatePasswordHandler(req *UpdatePasswordRequest) map[string]interface{} {
	log.Printf("Entering UpdatePasswordHandler")
	defer log.Printf("Exiting UpdatePasswordHandler")

	_, u, resp := EnsureAuthentication(req.Token)
	if len(resp) > 0 {
		return resp
	}

	minPasswordLength := 6
	if len(req.NewPassword) < minPasswordLength {
		log.Printf("New password too short")
		return map[string]interface{}{
			"error": fmt.Sprintf("New password must be at least %d characters long", minPasswordLength),
		}
	}

	u, err := user.Authenticate(u.Email, req.Password)
	if err != nil {
		log.Printf("User supplied incorrect current password")
		return map[string]interface{}{
			"error": "Incorrect current password",
		}
	}

	u.SetPassword(req.NewPassword)
	err = u.Insert()
	if err != nil {
		log.Printf("Error inserting updated user: %s", err.Error())
		return map[string]interface{}{
			"error": "Unable to update password",
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}
