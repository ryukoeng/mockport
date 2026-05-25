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
