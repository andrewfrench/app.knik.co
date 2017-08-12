package handlers

import (
	"log"
	"knik.co/library/account"
)

type GetAuthCodeRequest struct {
	Token string
	Username string
}

func GetAuthCodeHandler(req *GetAuthCodeRequest) map[string]interface{} {
	log.Printf("Entering GetAuthCodeHandler")
	defer log.Printf("Exiting GetAuthCodeHandler")

	_, u, resp := EnsureAuthentication(req.Token)
	if len(resp) > 0 {
		return resp
	}

	code := account.CodeGen(req.Username, u.Id)

	return map[string]interface{}{
		"auth_code": code,
	}
}
