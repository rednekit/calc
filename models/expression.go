package models

// Expression представляет арифметическое выражение.
type Expression struct {
    ID     int     `json:"id"`
    Status string  `json:"status"`
    Result float64 `json:"result"`
}
