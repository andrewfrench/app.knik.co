package handlers

import (
	"log"
	"knik.co/library/user"
)

type GetUsersRequest struct {
	Token string
}

func GetUsersHandler(req *GetUsersRequest) map[string]interface{} {
	log.Printf("Entering GetUsersHandler")
	defer log.Printf("Exiting GetUsersHandler")

	_, _, resp := EnsureAdminAuthentication(req.Token)
	if len(resp) > 0 {
		return resp
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
