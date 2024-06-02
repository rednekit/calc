package agent

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "your-module-name/models"
)

// TaskResult представляет результат выполнения задачи агентом.
type TaskResult struct {
    ID     int     `json:"id"`
    Result float64 `json:"result"`
}

// Task представляет задачу для вычисления арифметической операции.
type Task struct {
    ID           int    `json:"id"`
    Argument1    int    `json:"arg1"`
    Argument2    int    `json:"arg2"`
    Operation    string `json:"operation"`
    OperationTime int    `json:"operation_time"`
}

var httpClient = &http.Client{
    Timeout: 10 * time.Second,
}

func StartAgent(ctx context.Context, taskQueue chan<- models.Task, resultQueue <-chan models.TaskResult, endpoint string) {
    ticker := time.NewTicker(5 * time.Second) // Периодичность опроса оркестратора на предмет новых задач
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            task, err := getTaskFromOrchestrator(endpoint)
            if err != nil {
                fmt.Println("Failed to get task from orchestrator:", err)
                continue
            }

            if task != nil {
                taskQueue <- *task
            }

        case result := <-resultQueue:
            err := submitTaskResultToOrchestrator(endpoint, result)
            if err != nil {
                fmt.Println("Failed to submit task result to orchestrator:", err)
            }
        }
    }
}

func getTaskFromOrchestrator(endpoint string) (*models.Task, error) {
    resp, err := httpClient.Get(endpoint)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var task models.Task
    if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
        return nil, err
    }

    return &task, nil
}

func submitTaskResultToOrchestrator(endpoint string, result models.TaskResult) error {
    requestBody, err := json.Marshal(result)
    if err != nil {
        return err
    }

    resp, err := httpClient.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
}
