package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func saveTasks(tasks []Task) {
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error encoding JSON", err)
		os.Exit(1)
	}
	file.Close()

}

func readTasks() []Task {
	var tasks []Task
	jsonFile, err := os.Open("tasks.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		fmt.Println(err)
		os.Exit(1)
	}

	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&tasks)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
	return tasks
}

func addTask(description string) {
	tasks := readTasks()
	newTask := Task{
		ID:          len(tasks) + 1,
		Description: description,
	}

	tasks = append(tasks, newTask)

	saveTasks(tasks)

	fmt.Println("Task added")

}

func listTasks() {
	tasks := readTasks()
	if len(tasks) == 0 {
		fmt.Println("There are no tasks")
		return
	}
	for i := 0; i < len(tasks); i++ {
		status := " "
		if tasks[i].Done {
			status = "x"
		}
		fmt.Printf("%d: ["+status+"] %s\n", tasks[i].ID, tasks[i].Description)
	}

}

func deleteTasks(id int) {
	tasks := readTasks()
	newTasks := []Task{}

	for _, task := range tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)

		}
	}
	for i := range newTasks {
		newTasks[i].ID = i + 1
	}
	saveTasks(newTasks)
	fmt.Println("Task deleted succesfully")

}

func doneTask(id int) {
	tasks := readTasks()
	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}
	if !found {
		fmt.Println("task not found")
		os.Exit(1)
	}

	saveTasks(tasks)
	fmt.Println("Task marked as done")
}
func main() {

	if len(os.Args) < 2 {
		fmt.Println("You need to use a command")
		os.Exit(1)
	}
	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("you need to add a task")
			os.Exit(1)
		}
		description := os.Args[2]
		addTask(description)
	case "list":
		listTasks()
	case "delete":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			os.Exit(1)
		}

		deleteTasks(id)
	case "done":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			os.Exit(1)
		}
		doneTask(id)

	default:
		fmt.Println("not a valid command")
		os.Exit(1)
	}

}
