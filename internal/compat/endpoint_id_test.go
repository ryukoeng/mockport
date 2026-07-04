package compat

import (
	"net/http"
	"testing"
)

func TestEndpointIDNormalizesProviderPaths(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   string
	}{
		{http.MethodPost, "/slack/api/chat.postMessage", "post_slack_api_chat_postmessage"},
		{http.MethodGet, "/line/oauth2/v2.1/authorize", "get_line_oauth2_v2_1_authorize"},
		{http.MethodPost, "/line/v2/oauth/accessToken", "post_line_v2_oauth_accesstoken"},
		{http.MethodPost, "/line/v2/bot/chat/markAsRead", "post_line_v2_bot_chat_markasread"},
		{http.MethodGet, "/stripe/v1/checkout/sessions/{id}", "get_stripe_v1_checkout_sessions__id"},
		{http.MethodGet, "/openai/v1/responses/{id}", "get_openai_v1_responses__id"},
		{http.MethodGet, "/line/v2/bot/message/{messageId}/content", "get_line_v2_bot_message__messageid__content"},
		{http.MethodGet, "/line/v2/bot/profile/{userId}", "get_line_v2_bot_profile__userid"},
	}
	for _, tc := range tests {
		if got := endpointID(tc.method, tc.path); got != tc.want {
			t.Fatalf("endpointID(%q, %q) = %q, want %q", tc.method, tc.path, got, tc.want)
		}
	}
}
