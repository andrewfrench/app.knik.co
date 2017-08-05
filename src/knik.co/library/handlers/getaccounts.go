package handlers

import (
	"log"
	"knik.co/library/user"
	"knik.co/library/session"
)

type GetUsersRequest struct {
	Token string
}

func GetUsersHandler(req *GetUsersRequest) map[string]interface{} {
	log.Printf("Entering GetUsersHandler")
	defer log.Printf("Exiting GetUsersHandler")

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

	if !u.Admin {
		log.Printf("User is not admin")
		return map[string]interface{}{
			"error": "Unauthenticated",
		}
	}

	users, err := user.GetUsers()
	if err != nil {
		log.Printf("Error getting users: %s", err.Error())
		return map[string]interface{}{
			"error": "Error getting users",
		}
	}

	return map[string]interface{}{
		"users": users,
	}
}
