package line

type lineErrorDetail struct {
	Message  string `json:"message"`
	Property string `json:"property"`
}

type lineProfileResponse struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

type linePayResponse struct {
	ReturnCode    string      `json:"returnCode"`
	ReturnMessage string      `json:"returnMessage"`
	Info          linePayInfo `json:"info"`
}

type linePayInfo struct {
	TransactionID string `json:"transactionId"`
	OrderID       any    `json:"orderId"`
	Amount        any    `json:"amount"`
	Currency      any    `json:"currency"`
	Status        any    `json:"status"`
}

type richMenuData map[string]any
