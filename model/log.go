package model

type SwingLog struct {
	UserID   string  `json:"user_id"`
	Symbol   string  `json:"symbol"`
	Stoploss float64 `json:"stoploss"`
	Target   float64 `json:"target"`
	Account  Account `json:"account"`
}

type Account struct {
	Name     string  `json:"name"`
	UserID   string  `json:"user_id"`
	UserEmail string `json:"user_email"`
}
