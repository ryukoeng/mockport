package security

import (
	"net"
	"net/url"
	"strings"
)

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

func IsLoopbackRemoteAddr(remoteAddr string) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	return IsLoopbackHost(host)
}

func IsLoopbackHost(host string) bool {
	host = strings.Trim(strings.ToLower(strings.TrimSpace(host)), "[]")
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func IsSafeWebhookTargetURL(value string) bool {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" {
		return false
	}
	switch parsed.Scheme {
	case "http", "https":
	default:
		return false
	}
	if parsed.User != nil {
		return false
	}
	host := normalizedURLHost(parsed)
	switch host {
	case "localhost", "host.docker.internal", "app":
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func IsSafeOAuthRedirectURL(value string) bool {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" {
		return false
	}
	switch parsed.Scheme {
	case "http", "https":
	default:
		return false
	}
	if parsed.User != nil {
		return false
	}
	host := normalizedURLHost(parsed)
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func normalizedURLHost(parsed *url.URL) string {
	host := parsed.Hostname()
	if host == "" {
		host = parsed.Host
	}
	return strings.Trim(strings.ToLower(strings.TrimSpace(host)), "[]")
}
