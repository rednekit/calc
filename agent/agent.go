package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"calc/models"
)

var (
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	timeAdditionMs       int
	timeSubtractionMs    int
	timeMultiplicationMs int
	timeDivisionMs       int
)

func init() {
	var err error
	timeAdditionMs, err = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if err != nil {
		timeAdditionMs = 1000 // Default 1 second
	}
	timeSubtractionMs, err = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if err != nil {
		timeSubtractionMs = 1000 // Default 1 second
	}
	timeMultiplicationMs, err = strconv.Atoi(os.Getenv("TIME_MULTIPLICATION_MS"))
	if err != nil {
		timeMultiplicationMs = 1000 // Default 1 second
	}
	timeDivisionMs, err = strconv.Atoi(os.Getenv("TIME_DIVISION_MS"))
	if err != nil {
		timeDivisionMs = 1000 // Default 1 second
	}
}

func StartAgent(ctx context.Context, taskQueue chan models.Task, resultQueue chan models.TaskResult, endpoint string) {
	go fetchTasks(ctx, taskQueue, endpoint)
	go processTasks(ctx, taskQueue, resultQueue)
	go submitResults(ctx, resultQueue, endpoint)
}

func fetchTasks(ctx context.Context, taskQueue chan models.Task, endpoint string) {
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
			log.Println("Fetching task:", task)
			if task != nil {
				taskQueue <- *task
			}
		}
	}
}

func processTasks(ctx context.Context, taskQueue chan models.Task, resultQueue chan models.TaskResult) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-taskQueue:
			log.Println("Processing task:", task)
			result := processTask(task)
			resultQueue <- result
		}
	}
}

func submitResults(ctx context.Context, resultQueue chan models.TaskResult, endpoint string) {
	for {
		select {
		case <-ctx.Done():
			return
		case result := <-resultQueue:
			err := submitTaskResultToOrchestrator(endpoint, result)
			if err != nil {
				fmt.Println("Failed to submit task result to orchestrator:", err)
			}
		}
	}
}

func processTask(task models.Task) models.TaskResult {
	var result float64

	switch task.Operation {
	case "+":
		//time.Sleep(time.Duration(timeAdditionMs) * time.Millisecond)
		result = float64(task.Arg1) + float64(task.Arg2)
	case "-":
		//time.Sleep(time.Duration(timeSubtractionMs) * time.Millisecond)
		result = float64(task.Arg1) - float64(task.Arg2)
	case "*":
		//time.Sleep(time.Duration(timeMultiplicationMs) * time.Millisecond)
		result = float64(task.Arg1) * float64(task.Arg2)
	case "/":
		//time.Sleep(time.Duration(timeDivisionMs) * time.Millisecond)
		if task.Arg2 != 0 {
			result = float64(task.Arg1) / float64(task.Arg2)
		} else {
			result = 0 // handle division by zero appropriately
		}
	}

	return models.TaskResult{
		ID:     task.ID,
		Result: result,
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
