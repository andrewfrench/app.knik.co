package main

import (
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"log"
	"knik.co/library/handlers"
)

type incoming struct {
	Resource string `json:"resource"`
	Token string `json:"token"`
	Email string `json:"email"`
	Password string `json:"password"`
	AccountId string `json:"account_id"`
	Market string `json:"market"`
	Location string `json:"location"`
	Experience string `json:"experience"`
	Summary string `json:"summary"`
	Username string `json:"username"`
	AuthCode string `json:"auth_code"`
}

func Handle(evt *incoming, ctx *runtime.Context) (interface{}, error) {
	log.Printf("Entering Handle")
	defer log.Printf("Exiting Handle")

	switch evt.Resource {
	case "signin":
		return handlers.SignInHandler(&handlers.SignInRequest{
			Email: evt.Email,
			Password: evt.Password,
		}), nil

	case "signup":
		return handlers.SignUpHandler(&handlers.SignUpRequest{
			Email: evt.Email,
			Password: evt.Password,
		}), nil

	case "signout":
		return handlers.SignOutHandler(&handlers.SignOutRequest{
			Token: evt.Token,
		}), nil

	case "get-self":
		return handlers.GetSelfHandler(&handlers.GetSelfRequest{
			Token: evt.Token,
		}), nil

	case "get-account":
		return handlers.GetAccountHandler(&handlers.GetAccountRequest{
			Token: evt.Token,
			AccountId: evt.AccountId,
		}), nil

	case "update-account":
		return handlers.UpdateAccountHandler(&handlers.UpdateAccountRequest{
			Token: evt.Token,
			AccountId: evt.AccountId,
			Market: evt.Market,
			Location: evt.Location,
			Experience: evt.Experience,
			Summary: evt.Summary,
		}), nil

	case "get-auth-code":
		return handlers.GetAuthCodeHandler(&handlers.GetAuthCodeRequest{
			Token: evt.Token,
			Username: evt.Username,
		}), nil

	case "verify":
		return handlers.VerifyHandler(&handlers.VerifyRequest{
			Token: evt.Token,
			Username: evt.Username,
			AuthCode: evt.AuthCode,
		}), nil
	}

	return map[string]interface{}{
		"error": "Undefined resource",
	}, nil
}
