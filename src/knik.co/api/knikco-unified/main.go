package main

import (
	"log"
	"knik.co/library/handlers"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"strings"
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
	NewPassword string `json:"new_password"`
}

func Handle(evt *incoming, ctx *runtime.Context) (interface{}, error) {
	log.Printf("Entering Handle")
	defer log.Printf("Exiting Handle")

	evt.Username = strings.ToLower(evt.Username)
	evt.Email = strings.ToLower(evt.Email)

	switch evt.Resource {
	case "interest-email":
		return handlers.InterestEmailHander(&handlers.InterestEmailRequest{
			Email: evt.Email,
		}), nil

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

	case "get-users":
		return handlers.GetUsersHandler(&handlers.GetUsersRequest{
			Token: evt.Token,
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

	case "update-username":
		return handlers.UpdateUsernameHandler(&handlers.UpdateUsernameRequest{
			Token: evt.Token,
			AccountId: evt.AccountId,
			Username: evt.Username,
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

	case "update-password":
		return handlers.UpdatePasswordHandler(&handlers.UpdatePasswordRequest{
			Token: evt.Token,
			Password: evt.Password,
			NewPassword: evt.Password,
		}), nil
	}

	return map[string]interface{}{
		"error": "Undefined resource",
	}, nil
}
