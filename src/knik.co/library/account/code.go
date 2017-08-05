package account

import (
	"strings"
	"fmt"
	"knik.co/library/common"
	"regexp"
)

func CodeGen(username, userId string) string {
	username = strings.ToLower(username)
	username = strings.TrimSpace(username)
	hashInput := fmt.Sprintf("kn%s%sik", username, userId)
	hashOutput := common.Hash(hashInput)[:8] // Limit to 8 characters
	authCode := fmt.Sprintf("kn%sik", hashOutput)

	return authCode
}

func ExtractCodeCandidates(input string) []string {
	r, _ := regexp.Compile("kn[0-9a-f]+ik")

	return r.FindAllString(input, -1)
}
