package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

func main() {

	command := flag.String("command", "", "Enter a command to run")

	flag.Parse()

	for {
		fmt.Println("Please enter new command")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		*command = scanner.Text()
		runCommand(*command)
		fmt.Println()
	}

}

func runCommand(cmd string) {
	fmt.Println()
	switch cmd {
	case "add":
		addTask()
	case "delete":
		deleteTask()
	case "update":
		updateTask()
	case "update-status":
		updateTaskStatus()
	case "list":
		listTask()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Please enter valid command (add | update | update-status | delete | list | exit)")
	}
}

func getTasks() []Task {
	// Open file task.json if don't have then create it
	f, err := os.OpenFile("./task.json", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf(err.Error())
	}

	defer f.Close()

	decoder := json.NewDecoder(f)

	var tasks []Task

	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return tasks
}

func saveTask(task []Task) {

	f, err := os.OpenFile("./task.json", os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")

	err = encoder.Encode(task)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getId() int {
	tasks := getTasks()

	if len(tasks) == 0 {
		return 1
	}
	lastItem := tasks[len(tasks)-1]

	return lastItem.ID + 1
}

func isEmptyTask() bool {
	tasks := getTasks()
	if len(tasks) == 0 {
		fmt.Println("List task is empty, please add new task.")
		return true
	}
	return false
}

func addTask() {
	var description string
	fmt.Println("Please enter task description")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	description = scanner.Text()

	now := time.Now().Local()
	task := Task{
		ID:          getId(),
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks := []Task(getTasks())
	tasks = append(tasks, task)
	saveTask(tasks)

	fmt.Println("Task Created\n")
	showTask(task.ID)
}

func updateTaskStatusbyId(taskId int, status Status) ([]Task, error) {

	tasks := getTasks()
	for i, task := range tasks {
		if task.ID == taskId {
			tasks[i].Status = status
			now := time.Now().Local()
			tasks[i].UpdatedAt = now
			return tasks, nil
		}
	}
	return tasks, errors.New("can't find task\n")
}

func updateTaskById(taskId int, description string) ([]Task, error) {

	tasks := getTasks()
	for i := range tasks {
		if taskId == tasks[i].ID {
			tasks[i].Description = description
			now := time.Now().Local()
			tasks[i].UpdatedAt = now
			return tasks, nil
		}
	}
	return tasks, errors.New("Can't find task")
}

func updateTaskStatus() {
	if isEmptyTask() {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter task id\n")
	scanner.Scan()
	taskId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Must enter valid number.\n")
	}

	fmt.Println("Please enter new status (todo | in-progress | done)")
	scanner.Scan()
	rgx, err := regexp.MatchString(`\s*\b(todo|in-progress|done)\b\s*`, string(scanner.Bytes()))
	if err != nil {
		fmt.Println("Please enter valid status (todo | in-progress | done)")
		return
	}
	if !rgx {
		fmt.Println("Please enter valid status (todo | in-progress | done)")
	}
	newStatus := Status(scanner.Text())
	tasks, err := updateTaskStatusbyId(taskId, newStatus)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
	saveTask(tasks)
	showTask(taskId)
	return
}

func updateTask() {
	if isEmptyTask() {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter task id\n")
	scanner.Scan()

	taskId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("must enter valid number")

	}

	fmt.Println("Enter new description")
	scanner.Scan()

	newDescription := scanner.Text()
	task, err := updateTaskById(taskId, newDescription)
	if err != nil {
		log.Fatal(err)
	}
	saveTask(task)
	fmt.Println("New description has been updated\n")
	showTask(taskId)
	return

}

func deleteTaskById(taskId int) ([]Task, error) {

	tasks := getTasks()
	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return tasks, nil
		}
	}
	return tasks, errors.New("Task not found\n")
}

func deleteTask() {
	if isEmptyTask() {
		return
	}
	fmt.Println("Please enter the task id to delete")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	taskId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Must enter valid number")
		return
	}
	tasks, err := deleteTaskById(taskId)
	if err != nil {
		fmt.Println(err)
		return
	}
	saveTask(tasks)
	fmt.Println("Task has been deleted")

}

func listTask() {
	isEmptyTask()
	tasks := getTasks()
	for _, task := range tasks {
		fmt.Printf("ID: %d\n", task.ID)
		fmt.Printf("Description: %s\n", task.Description)
		fmt.Printf("Status: %s\n", task.Status)
		fmt.Printf("Created At: %s\n", task.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Updated At: %s\n", task.UpdatedAt.Format(time.RFC3339))
		fmt.Println("----------------------")

	}
}

func showTask(taskId int) {
	tasks := getTasks()
	for i, task := range tasks {
		if taskId == tasks[i].ID {
			fmt.Printf("ID: %d\n", task.ID)
			fmt.Printf("Description: %s\n", task.Description)
			fmt.Printf("Status: %s\n", task.Status)
			fmt.Printf("Created At: %s\n", task.CreatedAt.Format(time.RFC3339))
			fmt.Printf("Updated At: %s\n", task.UpdatedAt.Format(time.RFC3339))
			fmt.Println("----------------------")
		}
	}

}
