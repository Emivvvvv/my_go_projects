package TaskManager

import (
	"github.com/gofiber/fiber/v2"
	"simpleTaskManager/DBTaskController"
)

var db = DBTaskController.InitDB()

func AddTask(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	title := data["title"]
	description := data["description"]

	if len(title) != 0 || len(description) != 0 {
		newTask := DBTaskController.GenerateTaskToAdd(title, description, "Pending...")
		db.AddTask(newTask)

		return c.Status(200).SendString("user successfully added to the database")
	} else {
		return c.Status(400).SendString("title or description can not be empty")
	}
}

func GetTask(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id := data["id"]

	if len(id) != 0 {
		task := db.GetTask(DBTaskController.StringToObjectID(id))
		if task == nil {
			return c.Status(500).SendString("Couldn't find the task with id provided")
		}
		return c.JSON(task)
	} else {
		return c.Status(400).SendString("id can not be empty")
	}
}

func ViewTasks(c *fiber.Ctx) error {
	return c.JSON(db.GetAllTasks())
}

func MarkTask(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id := data["id"]

	if len(id) != 0 {
		task := db.GetTask(DBTaskController.StringToObjectID(id))
		if task == nil {
			return c.Status(500).SendString("Couldn't find the task with id provided")
		}
		db.UpdateTask(DBTaskController.StringToObjectID(id), DBTaskController.UpdateStatus("Completed!"))
		return c.Status(200).SendString("user with provided id successfully marked")
	} else {
		return c.Status(400).SendString("id can not be empty")
	}
}

func DeleteTask(c *fiber.Ctx) error {
	oldId := c.Params("id")
	id := oldId[1:]

	if len(id) != 0 {
		task := db.GetTask(DBTaskController.StringToObjectID(id))
		if task == nil {
			return c.Status(500).SendString("Couldn't find the task with id provided")
		}
		db.DeleteTask(DBTaskController.StringToObjectID(id))
		return c.Status(200).SendString("user with provided id successfully deleted")
	} else {
		return c.Status(400).SendString("id can not be empty")
	}
}

func FilterTasksByStatus(c *fiber.Ctx) error {
	mode := c.Params("mode")

	if len(mode) != 0 {
		var filteredTasks *[]DBTaskController.Task
		if mode == ":0" || mode == ":pending" {
			filteredTasks = db.FilterTasksByStatus("Pending...")
		} else if mode == ":1" || mode == ":completed" {
			filteredTasks = db.FilterTasksByStatus("Completed!")
		} else {
			return c.Status(500).SendString("mode can only be 0 or pending for Pending... tasks and 1 or completed for Completed! tasks")
		}
		return c.JSON(*filteredTasks)
	} else {
		return c.Status(400).SendString("mode can not be empty")
	}
}
