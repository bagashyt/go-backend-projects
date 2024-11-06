package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
}

func main() {
	var command string
	fmt.Println("Welcome to Expense Tracker\n")

	for {
		fmt.Println("Please enter new command (add | delete | list)\n")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command = scanner.Text()
		runCommand(command)
		fmt.Println()
	}
}

func runCommand(cmd string) {
	switch cmd {
	case "add":
		addExpense()
	case "delete":
		deleteExpense()
	case "list":
		listExpenses()
	}
}

func getExpenses() []Expense {
	f, err := os.OpenFile("./expense.json", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	decoder := json.NewDecoder(f)

	var tasks []Expense

	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return tasks
}

func getId() int {
	expenses := getExpenses()

	if len(expenses) == 0 {
		return 1
	}
	lastItem := expenses[len(expenses)-1]
	return lastItem.ID + 1
}

func saveExpense(expense []Expense) {
	f, err := os.OpenFile("./expense.json", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}

	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")

	err = encoder.Encode(expense)

	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}
}

func addExpense() {
	var description string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("input expense description")
	scanner.Scan()
	description = scanner.Text()

	fmt.Println("input expense amount")
	scanner.Scan()
	amount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("input valid number")
	}

	expense := Expense{
		ID:          getId(),
		Description: description,
		Date:        time.Now(),
		Amount:      amount,
	}
	expenses := []Expense(getExpenses())
	expenses = append(expenses, expense)
	saveExpense(expenses)

	fmt.Println("Expense Created")

}

func deleteExpenseById(id int) ([]Expense, error) {
	expenses := getExpenses()
	for i, expense := range expenses {
		if expense.ID == id {
			expenses = append(expenses[:i], expenses[i+1:]...)
			return expenses, nil
		}
	}
	return expenses, errors.New("expense not found")
}

func deleteExpense() {
	fmt.Println("please enter Expense ID")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	expenseId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("must enter valid number")
		return
	}
	expenses, err := deleteExpenseById(expenseId)
	if err != nil {
		fmt.Println(err)
		return
	}
	saveExpense(expenses)
	fmt.Println("expense has been deleted")
}

func listExpenses() {
	expenses := getExpenses()
	for _, expense := range expenses {
		fmt.Printf("ID: %d\n", expense.ID)
		fmt.Printf("Description: %s\n", expense.Description)
		fmt.Printf("Date: %s\n", expense.Date)
		fmt.Printf("Amount: %d $\n", expense.Amount)
		fmt.Println("----------------------")
	}
}
