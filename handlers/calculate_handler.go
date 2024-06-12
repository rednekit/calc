package handlers

import (
	"calc/models"
	"calc/storage"
	"calc/taskcreator"
	"calc/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CalculateHandler struct {
	Storage *storage.Storage
}

func (h *CalculateHandler) CalculateExpression(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid request payload")
		return
	}

	log.Println("Received expression:", req.Expression)

	tasks, err := taskcreator.CreateTasksFromExpression(req.Expression)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid expression")
		log.Println("Invalid expression")
		return
	}
	log.Println("Created tasks:", tasks)

	expression := models.Expression{
		Expression: req.Expression,
		Status:     "pending",
	}

	expressionID := h.Storage.AddExpression(expression)

	for _, task := range tasks {
		log.Println("Created taskID:", task.ID)
		log.Println("Created taskArg1:", task.Arg1)
		log.Println("Created taskArg2:", task.Arg2)
		log.Println("Created taskOperation:", task.Operation)
		log.Println("Created taskstatus:", task.Status)
		h.Storage.AddTask(task)

	}

	log.Println("Waiting to end multiply")
	time.Sleep(time.Second * 15)

	//	req2 := strings.ReplaceAll(req.Expression, " ", "")
	req2 := strings.TrimSpace(req.Expression)
	log.Println("Tasks :", tasks)
	for _, task := range tasks {
		log.Println("Task :", task)
		log.Println("Searching for taskID: ", task.ID)
		res, err := h.Storage.GetTaskResult(task.ID)
		log.Println("task id :", task.ID)
		log.Println("res :", res)
		if err != nil {
			panic(err)
		}
		i := strings.Index(req2, strconv.Itoa(task.Arg1)+" "+task.Operation+" "+strconv.Itoa(task.Arg2))
		req2 = req2[:i] + " " + strconv.FormatFloat(res.Result, 'f', -1, 32) + " " + req2[i+len(strconv.Itoa(task.Arg1)+" "+task.Operation+" "+strconv.Itoa(task.Arg2)):]
	}
	log.Println("req2 is:", req2)
	tasks, err = taskcreator.CreateTasksFromExpression(req2)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid expression")
		return
	}
	log.Println("Created tasks:", tasks)

	expression = models.Expression{
		Expression: req.Expression,
		Status:     "pending",
	}

	expressionID = h.Storage.AddExpression(expression)

	for _, task := range tasks {
		log.Println("Created taskID:", task.ID)
		log.Println("Created taskArg1:", task.Arg1)
		log.Println("Created taskArg2:", task.Arg2)
		log.Println("Created taskOperation:", task.Operation)
		log.Println("Created taskstatus:", task.Status)
		h.Storage.AddTask(task)

	}

	response := map[string]int{"id": expressionID}
	utils.RespondWithJSON(w, http.StatusCreated, response)
}
