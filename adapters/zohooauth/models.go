package zohooauth

// tokenResponse is the success body of the token exchange. The Zoho client only
// reads access_token.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

// errorResponse carries the error reason. Callers inspect the presence of the
// error field rather than the HTTP status code.
type errorResponse struct {
	Error string `json:"error"`
}

// userInfoResponse mirrors Zoho's user info shape. The keys are capitalized
// (Email / Display_Name); the client treats an empty Email as an error.
type userInfoResponse struct {
	Email       string `json:"Email"`
	DisplayName string `json:"Display_Name"`
}
