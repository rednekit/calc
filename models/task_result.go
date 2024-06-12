package models
// TaskResult представляет результат выполнения задачи агентом.
type TaskResult struct {
    ID     int     `json:"id"`
    Result float64 `json:"result"`
}
