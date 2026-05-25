package security

import "strings"

func RedactSecret(value string) string {
	return RedactValue(value)
}

func RedactValue(value string) string {
	if key, val, ok := strings.Cut(value, "="); ok {
		return key + "=" + RedactValue(val)
	}
	if LooksLikeSecret(value) {
		return "[real-looking secret redacted]"
	}
	if LooksLikeExternalServiceURL(value) {
		return "[external service URL redacted]"
	}
	if len(value) <= 12 {
		return "[redacted]"
	}
	return value[:9] + "..." + value[len(value)-4:]
}

func RedactMessage(message string) string {
	parts := strings.Fields(message)
	for i, part := range parts {
		trimmed := strings.Trim(part, ".,;:()[]{}\"'")
		if LooksLikeSecret(trimmed) || LooksLikeExternalServiceURL(trimmed) {
			parts[i] = strings.Replace(part, trimmed, RedactValue(trimmed), 1)
		}
	}
	return strings.Join(parts, " ")
}
