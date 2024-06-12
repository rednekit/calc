package main

import (
	"bytes"
	"calc/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	for {
		task, err := fetchTask()
		if err != nil {
			fmt.Println("No task available:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(1 * time.Second)
		if task.Status == "pending" {
			go func() int {
				task.Mu.Lock()
				fmt.Println("Goroutine is started")
				result := executeTask(task)
				fmt.Println("Result submited:", result)
				err = submitTaskResult(result)
				if err != nil {
					fmt.Println("Failed to submit task result:", err)
					task.Mu.Unlock()
					return 0
				}
				fmt.Println("Result for this id: ", result.ID, " is: ", result.Result)
				task.Mu.Unlock()
				return 0
			}()

		}

	}
}

func fetchTask() (models.Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")

	if err != nil {
		return models.Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Task{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Println("Fetching tasks:", resp)
	var response map[string]models.Task
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {

		return models.Task{}, err
	}
	log.Println("Fetching tasks:", response["task"])
	return response["task"], nil

}

func executeTask(task models.Task) models.TaskResult {
	var result float64
	log.Println("task.Operation:", task.Operation)
	switch task.Operation {
	case "+":
		result = float64(task.Arg1) + float64(task.Arg2)
		//time.Sleep(getOperationTime("TIME_ADDITION_MS"))
	case "-":
		result = float64(task.Arg1) - float64(task.Arg2)
		//time.Sleep(getOperationTime("TIME_SUBTRACTION_MS"))
	case "*":
		result = float64(task.Arg1) * float64(task.Arg2)
		//time.Sleep(getOperationTime("TIME_MULTIPLICATION_MS"))
	case "/":
		result = float64(task.Arg1) / float64(task.Arg2)
		//time.Sleep(getOperationTime("TIME_DIVISION_MS"))
	}
	log.Println("Result is:", result)
	return models.TaskResult{
		ID:     task.ID,
		Result: result,
	}
}

func submitTaskResult(result models.TaskResult) error {
	resultData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewReader(resultData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getOperationTime(envVar string) time.Duration {
	value, err := strconv.Atoi(os.Getenv(envVar))
	if err != nil {
		return 1 * time.Second
	}
	return time.Duration(value) * time.Millisecond
}
