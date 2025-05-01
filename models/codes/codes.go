package models

type CodeData struct {
	Codes []CodeBody `json:"codes"`
}

type CodeBody struct {
	Code string `json:"code"`
}

type Subscribers struct {
	Subscribers []Subscriber `json:"subscribers"`
}

type Subscriber struct {
	TGID   int64  `json:"TelegramID"`
	UserID string `json:"UserID"`
}
