package main

import (
	"github.com/gofiber/fiber/v2"
	"simpleTaskManager/TaskManager"
)

func Routers(app *fiber.App) {
	app.Get("/", hello)

	app.Post("/task", TaskManager.AddTask)
	app.Get("/task/:id", TaskManager.GetTask)
	app.Put("/task", TaskManager.MarkTask)
	app.Delete("/task/:id", TaskManager.DeleteTask)

	app.Get("/tasks", TaskManager.ViewTasks)
	app.Get("/tasks/:mode", TaskManager.FilterTasksByStatus)
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Welcome to Emivvvvv's dumbass uselles basicass task manager with fuckn database shit!")
}

func main() {
	app := fiber.New()
	Routers(app)
	app.Listen(":3000")
}
