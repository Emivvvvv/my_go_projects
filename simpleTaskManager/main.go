package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type task struct {
	id          int
	title       string
	description string
	completed   bool
}

var ID int = 0

func (t task) addTask(taskList *[]task) {
	*taskList = append(*taskList, t)
}

func (t task) printTask() {
	fmt.Println("Task ID: ", t.id)
	fmt.Println("Title: ", t.title)
	fmt.Println("Description: ", t.description)
	if t.completed {
		fmt.Println("Task is COMPLETED!\n")
	} else {
		fmt.Println("Task is not completed.\n")
	}
}

func (t task) markTask(taskList *[]task, i int) {
	if t.completed {
		fmt.Println("Task is already COMPLETED!")
	} else {
		(*taskList)[i].completed = true
		fmt.Println("Task is marked as COMPLETED!")
	}
}

func (t task) deleteTask(taskList []task, index int) []task {
	return append(taskList[:index], taskList[index+1:]...)
}

func generateTask(title string, description string) task {
	newTask := task{id: ID, title: title, description: description, completed: false}
	ID++
	return newTask
}

func main() {
	shouldContinue := true
	taskList := make([]task, 0)

	fmt.Print("\nWelcome to Emivvvvv's dumbass uselles basicass task manager!")

	for shouldContinue {

		fmt.Print("\n\nMenu:\n" +
			"1. Add a Task\n" +
			"2. View Tasks\n" +
			"3. Mark a Task as Completed\n" +
			"4. Delete a Task\n" +
			"5. Exit\n\n" +
			">>>")
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if err == nil {
			switch command {
			case "1":
				fmt.Print("Type the title of the task >>> ")
				title, err1 := reader.ReadString('\n')
				title = strings.TrimSpace(title)
				fmt.Print("\nType the description of the task >>> ")
				description, err2 := reader.ReadString('\n')
				description = strings.TrimSpace(description)

				if err1 != nil || err2 != nil {
					panic("Something went wrong at case 1")
				} else {
					newTask := generateTask(title, description)
					newTask.addTask(&taskList)
				}
				break

			case "2":
				if len(taskList) == 0 {
					fmt.Println("List is empty! You can add some task with typing command \"1\"")
				} else {
					for _, task := range taskList {
						task.printTask()
					}
				}
				break

			case "3":
				fmt.Print("Type the id of the task >>> ")
				idInput, err1 := reader.ReadString('\n')
				idInput = strings.TrimSpace(idInput)
				id, err1 := strconv.Atoi(idInput)

				if err1 == nil {
					for i, task := range taskList {
						if task.id == id {
							task.markTask(&taskList, i)
							break
						}
					}
					fmt.Println("Couldn't find the task with id: ", id)
				} else {
					panic("Something went wrong at case 3")
				}

			case "4":
				fmt.Print("Type the id of the task >>> ")
				idInput, err1 := reader.ReadString('\n')
				idInput = strings.TrimSpace(idInput)
				id, err1 := strconv.Atoi(idInput)

				if err1 == nil {
					for i, task := range taskList {
						if task.id == id {
							taskList = task.deleteTask(taskList, i)
							break
						}
					}
					fmt.Println("Couldn't find the task with id: ", id)
				} else {
					panic("Something went wrong at case 4")
				}

			case "5":
				shouldContinue = false
				fmt.Println("Exiting the program...")
				break
			}
		} else {
			panic(err)
		}
	}
}
