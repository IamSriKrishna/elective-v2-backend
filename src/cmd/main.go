package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/sk/elective/src/internal/config"
	"github.com/sk/elective/src/internal/delivery"
	"github.com/sk/elective/src/internal/repository"
	"github.com/sk/elective/src/internal/usecase"
	"github.com/sk/elective/src/pkg/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	bookingRepo := repository.NewCourseBookingRepository(db)

	// Initialize usecase
	authService := usecase.NewAuthService(studentRepo, cfg.JWT)
	courseService := usecase.NewCourseService(courseRepo, bookingRepo)

	// Initialize delivery
	authHandler := delivery.NewAuthHandler(authService)
	courseHandler := delivery.NewCourseHandler(courseService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // for testing, restrict in production
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	protected := api.Group("/", authHandler.AuthMiddleware)

	// Course routes
	courses := protected.Group("/courses")

	courses.Post("/", courseHandler.CreateCourse)
	courses.Get("/available", courseHandler.GetAvailableCourses)
	courses.Post("/book", courseHandler.BookCourse)
	courses.Get("/my-bookings", courseHandler.GetMyBookings)
	courses.Get("/all", courseHandler.GetAllCourses)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Course Booking API is running",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
