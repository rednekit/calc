package main

import (
	"calc/handlers"
	"calc/mux-main"
	"calc/storage"
	"log"
	"net/http"
)

func main() {
	store := storage.NewStorage()

	router := mux.NewRouter()
	calculateHandler := handlers.CalculateHandler{Storage: store}
	expressionHandler := handlers.ExpressionHandler{Storage: store}
	taskHandler := handlers.TaskHandler{Storage: store}

	router.HandleFunc("/api/v1/calculate", calculateHandler.CalculateExpression).Methods("POST")
	router.HandleFunc("/api/v1/expressions", expressionHandler.GetExpressions).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id:[0-9]+}", expressionHandler.GetExpressionByID).Methods("GET")
	router.HandleFunc("/internal/task", taskHandler.GetTask).Methods("GET")
	router.HandleFunc("/internal/task", taskHandler.SubmitTaskResult).Methods("POST")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
