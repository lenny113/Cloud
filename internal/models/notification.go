package models

type RegisterWebhook struct {
	Url                   string                 `json:"url"`
	Country               string                 `json:"country"`
	Event                 string                 `json:"event"`
	ThresholdNotification *ThresholdNotification `json:"threshold,omitempty" firestore:"threshold,omitempty"`
}

type ThresholdNotification struct {
	Field    string  `json:"field"`
	Operator string  `json:"operator"`
	Value    float64 `json:"value"`
}

type RegisteredWebhookResponse struct {
	Id      string `json:"id"`
	Country string `json:"country"`
	Event   string `json:"event"`
	Time    string `json:"time"`
}

type AllRegisteredWebhook struct {
	Id string `json:"id"`
	RegisterWebhook
}
