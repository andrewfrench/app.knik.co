package handlers

import (
	"knik.co/library/session"
	"knik.co/library/user"
	"log"
)

func EnsureAuthentication(token string) (*session.Session, *user.User, map[string]interface{}) {
	log.Printf("Entering EnsureAuthentication")
	defer log.Printf("Exiting EnsureAuthentication")

	unauthenticatedResponse := map[string]interface{}{
		"error": "Unauthenticated",
		"unauthenticated": true,
	}

	s, err := session.GetSessionBySessionId(token)
	if err != nil {
		return &session.Session{}, &user.User{}, unauthenticatedResponse
	}

	u, err := user.GetUserById(s.UserId)
	if err != nil {
		return &session.Session{}, &user.User{}, unauthenticatedResponse
	}

	return s, u, map[string]interface{}{}
}

func EnsureAdminAuthentication(token string) (*session.Session, *user.User, map[string]interface{}) {
	log.Printf("Entering EnsureAdminAuthentication")
	defer log.Printf("Exiting EnsureAdminAuthentication")

	bouncedResponse := map[string]interface{}{
		"error": "Insufficient privileges",
		"bounce": true,
	}

	s, u, unauthenticatedResponse := EnsureAdminAuthentication(token)
	if len(unauthenticatedResponse) > 0 {
		return s, u, unauthenticatedResponse
	}

	if !u.Admin {
		return s, u, bouncedResponse
	}

	return s, u, map[string]interface{}{}
}
