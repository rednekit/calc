package storage

import (
	"calc/models"
	"errors"
	"log"
	"sync"
)

var nextTaskID int = 1

type Storage struct {
	expressions map[int]models.Expression
	tasks       map[int]models.Task
	results     map[int]models.TaskResult
	mutex       sync.Mutex
	nextID      int
}

func NewStorage() *Storage {
	return &Storage{
		expressions: make(map[int]models.Expression),
		tasks:       make(map[int]models.Task),
		results:     make(map[int]models.TaskResult),
		nextID:      1,
	}
}

func (s *Storage) AddExpression(expression models.Expression) int {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()

	id := s.nextID
	s.nextID++
	expression.ID = id
	s.expressions[id] = expression

	return id
}

func (s *Storage) GetExpressions() []models.Expression {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()

	expressions := make([]models.Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		expressions = append(expressions, expr)
	}

	return expressions
}

func (s *Storage) GetExpressionByID(id int) (models.Expression, error) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()

	expr, exists := s.expressions[id]
	if !exists {
		return models.Expression{}, errors.New("expression not found")
	}

	return expr, nil
}

func (s *Storage) AddTask(task models.Task) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()

	task.ID = nextTaskID
	nextTaskID++
	log.Println("Saving task:", task)
	s.tasks[task.ID] = task
}

func (s *Storage) GetTask() (models.Task, error) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()

	for _, task := range s.tasks {
		return task, nil
	}

	return models.Task{}, errors.New("no tasks available")
}

func (s *Storage) SubmitTaskResult(result models.TaskResult) error {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()

	_, exists := s.tasks[result.ID]
	if !exists {
		return errors.New("task not found")
	}
	log.Println("Task Result saving :", result)
	s.results[result.ID] = result
	log.Println("Task Result saving :", s.results[result.ID])
	var x = s.tasks[result.ID]
	x.Status = "done"
	s.tasks[result.ID] = x
	//delete(s.tasks, result.ID)

	return nil
}

func (s *Storage) GetTaskResult(id int) (models.TaskResult, error) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()

	result, exists := s.results[id]
	if !exists {
		return models.TaskResult{}, errors.New("result not found")
	}

	return result, nil
}
