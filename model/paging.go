package model

type Page struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	// SortType    string `json:"sort_type"`
	// SortParam   string `json:"sort_param"`
	// Total       uint   `json:"total"`
}
