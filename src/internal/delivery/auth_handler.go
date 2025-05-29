package delivery

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sk/elective/src/internal/domain"
)

type AuthHandler struct {
    authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
    RegisterNo string `json:"register_no" validate:"required"`
    Password   string `json:"password" validate:"required,min=6"`
    Department string `json:"department"`
}

type LoginRequest struct {
    RegisterNo string `json:"register_no" validate:"required"`
    Password   string `json:"password" validate:"required"`
}

type AuthResponse struct {
    Token   string      `json:"token"`
    Student interface{} `json:"student"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if req.Department == "" {
        req.Department = "CSE"
    }

    student, err := h.authService.Register(req.RegisterNo, req.Password, req.Department)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Student registered successfully",
        "student": fiber.Map{
            "id":          student.ID,
            "register_no": student.RegisterNo,
            "department":  student.Department,
        },
    })
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    token, student, err := h.authService.Login(req.RegisterNo, req.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(AuthResponse{
        Token: token,
        Student: fiber.Map{
            "id":          student.ID,
            "register_no": student.RegisterNo,
            "department":  student.Department,
        },
    })
}

func (h *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Authorization header required",
        })
    }

    tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
    student, err := h.authService.ValidateToken(tokenString)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    c.Locals("student", student)
    return c.Next()
}

