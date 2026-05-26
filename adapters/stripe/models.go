package stripe

type errorBody struct {
	Error stripeError `json:"error"`
}

type stripeError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message"`
}
