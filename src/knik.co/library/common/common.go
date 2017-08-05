package common

import (
	"fmt"
	"crypto/sha256"
)

func Hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}
