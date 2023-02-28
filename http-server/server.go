package httpserver

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}

type Server struct {
	Router *fiber.App
}

func MakeServer() Server {
	return Server{
		Router: fiber.New(),
	}
}

func (s *Server) Start() {
	allowedOrigins := os.Getenv("ALLOW_ORIGINS")
	redisUrl := os.Getenv("REDIS_URL")

	var limiterConfig limiter.Config
	if redisUrl != "" {
		redisStore := redis.New(redis.Config{
			URL:   redisUrl,
			Reset: false,
		})

		limiterConfig = limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.IP() == "127.0.0.1"
			},
			Max:        30,
			Expiration: 60 * time.Second,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(Response[any]{
					Status:  -1,
					Message: "Rate limited",
					Data:    nil,
				})
			},
			Storage: redisStore,
		}
	} else {
		limiterConfig = limiter.ConfigDefault
	}

	s.Router.Use(limiter.New(limiterConfig))

	if allowedOrigins == "" {
		s.Router.Use(cors.New())
	} else {
		s.Router.Use(cors.New(cors.Config{
			AllowOrigins:     allowedOrigins,
			AllowCredentials: true,
		}))
	}

	s.Router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	s.Router.Post("/login", Login)
	s.Router.Get("/me", Me)
	s.Router.Get("/home", Home)
	s.Router.Get("/photo", UserPhoto)

	academicGroup := s.Router.Group("/academic")
	academicGroup.Get("/summary", AcademicSummary)
	academicGroup.Get("/history", AcademicHistory)
	academicGroup.Get("/classes", CourseClasses)

	s.Router.Listen(getPort())
}
