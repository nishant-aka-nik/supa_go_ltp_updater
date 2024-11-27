package model

type SwingLog struct {
	ID           int     `json:"id"`
	UserID       string  `json:"user_id"`
	Symbol       string  `json:"symbol"`
	Stoploss     float64 `json:"stoploss"`
	Target       float64 `json:"target"`
	Pivot        float64 `json:"pivot"`
	BuyPrice     float64 `json:"buy_price"`
	LongPosition bool    `json:"long_position"`
	Account      Account `json:"account"`
}

type Account struct {
	Name           string `json:"name"`
	UserID         string `json:"user_id"`
	UserEmail      string `json:"user_email"`
	SecondaryEmail string `json:"secondary_email"`
}
