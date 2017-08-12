package handlers

import (
	"log"
)

type GetSelfRequest struct {
	Token string
}

func GetSelfHandler(req *GetSelfRequest) map[string]interface{} {
	log.Printf("Entering GetSelfHandler")
	defer log.Printf("Exiting GetSelfHandler")

	_, u, resp := EnsureAuthentication(req.Token)
	if len(resp) > 0 {
		return resp
	}

	return map[string]interface{}{
		"profile": u,
	}
}
