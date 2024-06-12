package handlers

import (
	"calc/models"
	"calc/storage"
	"calc/utils"
	"encoding/json"
	"log"
	"net/http"
	//"calc/mux-main"
)

type TaskHandler struct {
	Storage *storage.Storage
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.Storage.GetTask()
	log.Println("Storage task:", task)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "No tasks available")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]models.Task{"task": task})
}

func (h *TaskHandler) SubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var result models.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.Storage.SubmitTaskResult(result); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	w.WriteHeader(http.StatusOK)
}
