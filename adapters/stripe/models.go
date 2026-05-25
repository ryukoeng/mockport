package stripe

type errorBody struct {
	Error stripeError `json:"error"`
}

type stripeError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
