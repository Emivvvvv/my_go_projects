package main

import (
	"github.com/gofiber/fiber/v2"
	"simpleTaskManager/TaskManager"
)

func Routers(app *fiber.App) {
	app.Get("/", hello)
	app.Post("/user", TaskManager.AddTask)
	app.Get("/user:id", TaskManager.GetTask)
	app.Get("/users", TaskManager.ViewTasks)
	app.Put("/user", TaskManager.MarkTask)
	app.Delete("/user:id", TaskManager.DeleteTask)
	app.Get("/users:mode", TaskManager.FilterTasksByStatus)
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Welcome to Emivvvvv's dumbass uselles basicass task manager with fuckn database shit!")
}

func main() {
	app := fiber.New()
	Routers(app)
	app.Listen(":3000")
}
