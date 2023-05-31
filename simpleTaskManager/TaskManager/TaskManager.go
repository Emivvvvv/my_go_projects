package TaskManager

import (
	"github.com/gofiber/fiber/v2"
	"simpleTaskManager/DBTaskController"
)

var db = DBTaskController.InitDB()

type AddTaskBody struct {
	Title       string `json:"title" xml:"title" form:"title"`
	Description string `json:"description" xml:"description" form:"description"`
}

type MarkTaskBody struct {
	Id string `json:"id" xml:"id" form:"id"`
}

func AddTask(c *fiber.Ctx) error {
	t := new(AddTaskBody)
	if err := c.BodyParser(t); err != nil {
		return err
	}

	if len(t.Title) != 0 || len(t.Description) != 0 {
		newTask := DBTaskController.GenerateTaskToAdd(t.Title, t.Description, "Pending...")
		err := db.AddTask(newTask)

		if err != nil {
			return c.Status(200).SendString("user successfully added to the database")
		} else {
			return err
		}
	} else {
		return c.Status(400).SendString("title or description can not be empty")
	}
}

func GetTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if len(id) != 0 {
		id, err := DBTaskController.StringToObjectID(id)
		if err != nil {
			return err
		}

		task, err := db.GetTask(id)
		if task == nil {
			return err
		}

		return c.JSON(task)
	} else {
		return c.Status(400).SendString("id can not be empty")
	}
}

func ViewTasks(c *fiber.Ctx) error {
	tasks, err := db.GetAllTasks()

	if err != nil {
		return err
	}

	return c.JSON(tasks)
}

func MarkTask(c *fiber.Ctx) error {
	t := new(MarkTaskBody)
	if err := c.BodyParser(t); err != nil {
		return err
	}

	if len(t.Id) != 0 {
		id, err := DBTaskController.StringToObjectID(t.Id)
		if err != nil {
			return err
		}

		task, err := db.GetTask(id)
		if task == nil {
			return err
		}

		err = db.UpdateTask(id, DBTaskController.UpdateStatus("Completed!"))
		if err != nil {
			return err
		}

		return c.Status(200).SendString("user with provided id successfully marked")
	} else {
		return c.Status(400).SendString("id can not be empty")
	}
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if len(id) != 0 {
		id, err := DBTaskController.StringToObjectID(id)
		if err != nil {
			return err
		}

		task, err := db.GetTask(id)
		if task == nil {
			return err
		}

		err = db.DeleteTask(id)
		if err != nil {
			return err
		}

		return c.Status(200).SendString("user with provided id successfully deleted")
	} else {
		return c.Status(400).SendString("id can not be empty")
	}
}

func FilterTasksByStatus(c *fiber.Ctx) error {
	mode := c.Params("mode")

	if len(mode) != 0 {
		var filteredTasks *[]DBTaskController.Task
		var err error
		if mode == "0" || mode == "pending" {
			filteredTasks, err = db.FilterTasksByStatus("Pending...")
		} else if mode == "1" || mode == "completed" {
			filteredTasks, err = db.FilterTasksByStatus("Completed!")
		} else {
			return c.Status(500).SendString(mode)
			//return c.Status(500).SendString("mode can only be 0 or pending for Pending... tasks and 1 or completed for Completed! tasks")
		}
		if err != nil {
			return err
		}
		return c.JSON(*filteredTasks)
	} else {
		return c.Status(400).SendString("mode can not be empty")
	}
}
