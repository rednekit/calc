package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "your-module-name/models"
    "your-module-name/storage"
)

// GetExpressions возвращает список арифметических выражений.
func GetExpressions(w http.ResponseWriter, r *http.Request) {
    expressions := storage.GetExpressions() // Получаем список выражений из хранилища

    // Формируем ответ
    response := map[string]interface{}{"expressions": expressions}
    json.NewEncoder(w).Encode(response)
}

// GetExpressionByID возвращает арифметическое выражение по его идентификатору.
func GetExpressionByID(w http.ResponseWriter, r *http.Request) {
    // Получаем идентификатор выражения из URL
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid expression ID", http.StatusBadRequest)
        return
    }

    // Получаем выражение по его идентификатору из хранилища
    expression, err := storage.GetExpressionByID(id)
    if err != nil {
        if errors.Is(err, storage.ErrExpressionNotFound) {
            http.Error(w, "Expression not found", http.StatusNotFound)
        } else {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    // Формируем ответ
    response := map[string]interface{}{"expression": expression}
    json.NewEncoder(w).Encode(response)
}
