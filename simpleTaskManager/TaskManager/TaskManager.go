package TaskManager

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"simpleTaskManager/Auth"
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

func GetUserID(c *fiber.Ctx) (string, error) {
	sess, sessErr := Auth.Store.Get(c)
	if sessErr != nil {
		return "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong: " + sessErr.Error(),
		})
	}

	userID := sess.Get("USER_ID")
	if userID != nil {
		userIDStr, ok := userID.(string)
		if ok {
			return userIDStr, nil
		} else {
			return "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "user ID in session is not a string",
			})
		}
	} else {
		return "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "user ID not found in session",
		})
	}

}

func AddTask(c *fiber.Ctx) error {
	t := new(AddTaskBody)
	if err := c.BodyParser(t); err != nil {
		return err
	}

	if len(t.Title) != 0 || len(t.Description) != 0 {
		userID, idErr := GetUserID(c)
		if idErr != nil {
			return idErr
		}
		newTask := DBTaskController.GenerateTaskToAdd(userID, t.Title, t.Description, "Pending...")
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

		userID, idErr := GetUserID(c)
		if idErr != nil {
			return idErr
		}

		if task.TaskOwner != userID {
			return fmt.Errorf("user don't have the authority to update this task")
		}

		err = db.UpdateTask(id, DBTaskController.UpdateStatus("Completed!"))
		if err != nil {
			return err
		}

		return c.Status(200).SendString("task with provided id successfully marked")
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

		userID, idErr := GetUserID(c)
		if idErr != nil {
			return idErr
		}

		if task.TaskOwner != userID {
			return fmt.Errorf("user don't have the authority to delete this task")
		}

		err = db.DeleteTask(id)
		if err != nil {
			return err
		}

		return c.Status(200).SendString("task with provided id successfully deleted")
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
