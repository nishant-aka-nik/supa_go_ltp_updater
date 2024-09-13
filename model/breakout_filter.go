package model

type BreakoutFilter struct {
	ID              int     `json:"id"`
	Symbol          string  `json:"symbol"`
	CreatedAt       string  `json:"created_at"`
	CrossMatch      bool    `json:"cross_match"`
	CrossMatchPivot float64 `json:"cross_match_pivot"`
	CrossMatchDate  string  `json:"cross_match_date"`
	Entry           bool    `json:"entry"`
	EntryPrice      float64 `json:"entry_price"`
	EntryDate       string  `json:"entry_date"`
	VolumeTimes     float64 `json:"volume_times"`
}

func (s BreakoutFilter) GetEntry() float64 {
	return s.CrossMatchPivot + (s.CrossMatchPivot * 0.01)
}
