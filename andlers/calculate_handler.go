package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "your-module-name/models"
    "your-module-name/storage"
)

// CalculateExpression обрабатывает запрос на вычисление арифметического выражения.
func CalculateExpression(w http.ResponseWriter, r *http.Request) {
    var exp models.Expression
    if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Создаем задачи для агентов на основе выражения
    tasks := createTasksFromExpression(exp)

    // Отправляем задачи на выполнение агентам
    for _, task := range tasks {
        agentQueue <- task
    }

    // Отправляем ответ клиенту
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{"id": exp.ID})
}

// createTasksFromExpression создает задачи для агентов на основе арифметического выражения.
func createTasksFromExpression(exp models.Expression) []models.Task {
    // Здесь нужно добавить код для разбора выражения и создания задач для агентов
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
