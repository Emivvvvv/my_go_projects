package main

import (
	"bufio"
	"fmt"
	"os"
	"simpleTaskManager/TaskManager"
	"strings"
)

func input(info string) (string, error) {
	fmt.Print(info)
	reader := bufio.NewReader(os.Stdin)
	inputString, err := reader.ReadString('\n')
	inputString = strings.TrimSpace(inputString)
	fmt.Println()
	return inputString, err
}

func main() {
	shouldContinue := true
	db := TaskManager.InitDB()

	fmt.Print("\nWelcome to Emivvvvv's dumbass uselles basicass task manager with fuckn database shit!")

	menu := "\n\nMenu:\n" +
		"1. Add a Task\n" +
		"2. View Tasks\n" +
		"3. Mark a Task as Completed\n" +
		"4. Delete a Task\n" +
		"5. Filter tasks by status\n" +
		"6. Exit\n\n" +
		">>>"

	for shouldContinue {
		command, err := input(menu)

		if err == nil {
			switch command {
			case "1":
				title, err1 := input("Type the title of the task >>> ")
				description, err2 := input("Type the description of the task >>> ")

				if err1 == nil && err2 == nil {
					newTask := TaskManager.GenerateTaskToAdd(title, description, "Pending...")
					db.AddTask(newTask)
				} else {
					panic("Something went wrong at case 1")
				}
				break

			case "2":
				for _, task := range db.GetAllTasks() {
					task.Print()
				}
				break

			case "3":
				id, err1 := input("Type the id of the task >>> ")

				if err1 == nil {
					task := db.GetTask(TaskManager.StringToObjectID(id))
					if task == nil {
						fmt.Println("Couldn't find the task with id: ", id)
					}
					db.UpdateTask(TaskManager.StringToObjectID(id), TaskManager.UpdateStatus("Completed!"))
					break
				} else {
					panic("Something went wrong at case 3")
				}

			case "4":
				id, err1 := input("Type the id of the task >>> ")

				if err1 == nil {
					task := db.GetTask(TaskManager.StringToObjectID(id))
					if task == nil {
						fmt.Println("Couldn't find the task with id: ", id)
					}
					db.DeleteTask(TaskManager.StringToObjectID(id))
				} else {
					panic("Something went wrong at case 4")
				}

			case "5":
				userInput, err := input("Type\n" +
					"0 to filter pending tasks\n" +
					"1 to filter completed tasks\n" +
					">>>")
				if err == nil {
					var filteredTasks *[]TaskManager.Task
					if userInput == "0" {
						filteredTasks = db.FilterTasksByStatus("Pending...")
					} else if userInput == "1" {
						filteredTasks = db.FilterTasksByStatus("Completed!")
					} else {
						fmt.Println("Invalid input!")
					}
					for _, task := range *filteredTasks {
						task.Print()
					}
				} else {
					panic("Something went wrong at case 5")
				}
				break
			case "6":
				shouldContinue = false
				fmt.Println("Exiting the program...")
				break
			}
		} else {
			panic(err)
		}
	}
}
