package security

import "strings"

var dangerousURLPrefixes = []string{
	"https://api.stripe.com",
	"https://api.openai.com",
	"https://api.github.com",
	"https://api.line.me",
	"https://slack.com/api",
}

func LooksLikeExternalServiceURL(value string) bool {
	for _, prefix := range dangerousURLPrefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}
