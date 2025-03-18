package models

type CodeData struct {
	Codes []CodeBody `json:"codes"`
}

type CodeBody struct {
	Code string `json:"code"`
}

type Subscribers struct {
	Subscriber []int64 `json:"subscribers"`
}
