package httpserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	Router *fiber.App
}

func MakeServer() Server {
	return Server{
		Router: fiber.New(),
	}
}

func (s *Server) Start() {
	s.Router.Use(cors.New())
	s.Router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	s.Router.Post("/login", Login)
	s.Router.Get("/home", Home)
	s.Router.Get("/photo", UserPhoto)

	academicGroup := s.Router.Group("/academic")
	academicGroup.Get("/summary", AcademicSummary)
	academicGroup.Get("/history", AcademicHistory)

	s.Router.Listen(":8080")
}
