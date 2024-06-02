package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "your-module-name/handlers"
)

func main() {
    r := mux.NewRouter()
    
    // Handlers
    r.HandleFunc("/api/v1/calculate", handlers.CalculateExpression).Methods("POST")
    r.HandleFunc("/api/v1/expressions", handlers.GetExpressions).Methods("GET")
    r.HandleFunc("/api/v1/expressions/{id}", handlers.GetExpressionByID).Methods("GET")
    r.HandleFunc("/internal/task", handlers.GetTask).Methods("GET")
    r.HandleFunc("/internal/task", handlers.SubmitTaskResult).Methods("POST")

    // Run server
    log.Fatal(http.ListenAndServe(":8080", r))
}
