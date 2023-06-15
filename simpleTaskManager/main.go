package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"simpleTaskManager/Auth"
	"simpleTaskManager/TaskManager"
	"time"
)

func Routers(app *fiber.App) {
	app.Get("/", hello)

	app.Post("/auth/register", Auth.Register)
	app.Post("/auth/login", Auth.Login)
	app.Post("/auth/logout", Auth.Logout)
	app.Post("/auth/healthcheck", Auth.HealthCheck)

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
	store := session.New(session.Config{
		CookieHTTPOnly: true,
		//CookieSecure: true, for https
		Expiration: time.Hour * 5,
	})

	Auth.InitStore(store)

	app.Use(Auth.NewMiddleware(), cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
		AllowHeaders:     "Access-Control-Allow-Origin, Content-Type, Origin, Accept",
	}))
	Routers(app)
	app.Listen(":3000")
}
