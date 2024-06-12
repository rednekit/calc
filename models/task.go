package models

import "sync"

// Task представляет задачу для вычисления арифметической операции.
type Task struct {
	ID        int    `json:"id"`
	Arg1      int    `json:"arg1"`
	Arg2      int    `json:"arg2"`
	Operation string `json:"operation"`
	Status    string `json:"status"`
	Flag      int    `json:"flag"`
	Mu        sync.Mutex
}
