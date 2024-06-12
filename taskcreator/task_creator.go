package taskcreator

import (
	"calc/models"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func CreateTasksFromExpression(expression string) ([]models.Task, error) {
	flag := 0
	tokens := strings.Fields(expression)
	fmt.Println("Token: ", flag, tokens)
	if len(tokens) == 0 {
		return nil, errors.New("empty expression")
	}

	var tasks []models.Task
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}
	var stack []string
	var postfix []string

	for _, token := range tokens {
		fmt.Println("Token: ", token)
		switch token {
		case "+", "-", "*", "/":
			if precedence[token] == 2 && flag < 2 {
				flag = 2
			} else if precedence[token] == 1 && flag < 1 {
				flag = 1
			}
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {

				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		default:
			postfix = append(postfix, token)
		}
	}

	for len(stack) > 0 {
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	var argStack []float64
	i := 1
	for _, token := range postfix {
		fmt.Println("flag: ", flag)
		if flag == 2 {
			fmt.Println("Postfix Token: ", token, "Flag :", flag)
			switch token {
			case "+", "-":
				fmt.Println("+- ")
			case "*", "/":
				fmt.Println("case in")
				if len(argStack) < 2 {
					return nil, errors.New("invalid expression")
				}

				arg2 := argStack[len(argStack)-1]
				arg1 := argStack[len(argStack)-2]

				fmt.Println("Argstack: ", argStack)
				argStack = argStack[:len(argStack)-2]
				fmt.Println("Argstack: ", argStack)

				fmt.Println("Args: ", arg2, "Arg1 :", arg1)
				task := models.Task{
					ID:        i,
					Arg1:      int(arg1),
					Arg2:      int(arg2),
					Operation: token,
					Status:    "pending",
					Flag:      2,
				}
				i++
				fmt.Println("Task: ", task)
				tasks = append(tasks, task)

				//argStack = append(argStack, float64(task.ID))
			default:
				if token != "*" && token != "/" && token != "+" && token != "-" {
					value, err := strconv.ParseFloat(token, 64)
					if err != nil {
						return nil, err
					}
					argStack = append(argStack, value)
					fmt.Println("Argstack: ", argStack)
				}
			}
		} else if flag == 1 {
			fmt.Println("Postfix Token: ", token, "Flag :", flag)
			switch token {
			case "*", "/":
				fmt.Println("*/ ")
			case "+", "-":
				if len(argStack) < 2 {
					return nil, errors.New("invalid expression")
				}

				arg2 := argStack[len(argStack)-1]
				arg1 := argStack[len(argStack)-2]
				fmt.Println("Argstack: ", argStack)

				argStack = argStack[:len(argStack)-2]
				fmt.Println("Argstack: ", argStack)

				task := models.Task{
					ID:        i,
					Arg1:      int(arg1),
					Arg2:      int(arg2),
					Operation: token,
					Status:    "pending",
					Flag:      1,
				}
				i++

				tasks = append(tasks, task)

				//argStack = append(argStack, float64(task.ID))
			default:
				if token != "*" && token != "/" && token != "+" && token != "-" {
					value, err := strconv.ParseFloat(token, 64)
					if err != nil {
						return nil, err
					}
					argStack = append(argStack, value)
					fmt.Println("Argstack: ", argStack)
				}
			}
		}
	}
	fmt.Println("Argstack: ", argStack)

	// if len(argStack) != 1 {
	// 	return nil, errors.New("invalid expression")
	// }
	log.Println("Generated tasks:", tasks)

	return tasks, nil
}
