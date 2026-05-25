package security

import "strings"

var dangerousSecretPrefixes = []string{
	"sk_live_",
	"sk_test_",
	"AKIA",
	"ASIA",
	"ghp_",
	"github_pat_",
	"xoxb-",
	"xoxp-",
	"AIza",
	"whsec_",
}

var fakeSecretPrefixes = []string{
	"mockport_",
	"local_",
	"fake_",
	"dummy_",
}

var dangerousURLPrefixes = []string{
	"https://api.stripe.com",
	"https://api.openai.com",
	"https://api.github.com",
	"https://api.line.me",
	"https://slack.com/api",
}

func LooksLikeSecret(value string) bool {
	if value == "whsec_mockport" {
		return false
	}
	for _, prefix := range fakeSecretPrefixes {
		if strings.HasPrefix(value, prefix) {
			return false
		}
	}
	for _, prefix := range dangerousSecretPrefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func LooksLikeExternalServiceURL(value string) bool {
	for _, prefix := range dangerousURLPrefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func RedactSecret(value string) string {
	if LooksLikeSecret(value) {
		return "[real-looking secret redacted]"
	}
	if len(value) <= 12 {
		return "[redacted]"
	}
	return value[:9] + "..." + value[len(value)-4:]
}
