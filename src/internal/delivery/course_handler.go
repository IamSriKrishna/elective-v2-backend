package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sk/elective/src/internal/domain"
	"github.com/sk/elective/src/internal/repository/models"
)

type CourseHandler struct {
	courseService domain.CourseService
}

func NewCourseHandler(courseService domain.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

type CreateCourseRequest struct {
	Name           string   `json:"name" validate:"required"`
	PDFLink        string   `json:"pdf_link"`
	Rating         float64  `json:"rating"`
	StaffNames     []string `json:"staff_names"`
	ImageLink      string   `json:"image_link"`
	Description    string   `json:"description"`
	AvailableSeats int      `json:"available_seats"`
	Departments    []string `json:"departments" validate:"required"`
	Genres         []string `json:"genres"`
	CourseType     int      `json:"course_type" validate:"required,min=1,max=4"`
	TotalSeats     int      `json:"total_seats" validate:"required,min=1"`
}

type CourseResponse struct {
	ID             uint     `json:"id"`
	Name           string   `json:"name"`
	PDFLink        string   `json:"pdf_link"`
	Rating         float64  `json:"rating"`
	StaffNames     []string `json:"staff_names"`
	ImageLink      string   `json:"image_link"`
	Description    string   `json:"description"`
	Departments    []string `json:"departments"`
	Genres         []string `json:"genres"`
	CourseType     int      `json:"course_type"`
	TotalSeats     int      `json:"total_seats"`
	SeatsBooked    []string `json:"seats_booked"`
	AvailableSeats int      `json:"available_seats"`
}

type BookCourseRequest struct {
	CourseID uint   `json:"course_id" validate:"required"`
	SeatNo   string `json:"seat_no" validate:"required"`
}

func toCourseResponse(course models.Course) CourseResponse {
	return CourseResponse{
		ID:             course.ID,
		Name:           course.Name,
		PDFLink:        course.PDFLink,
		Rating:         course.Rating,
		StaffNames:     []string(course.StaffNames),
		ImageLink:      course.ImageLink,
		Description:    course.Description,
		Departments:    []string(course.Departments),
		Genres:         []string(course.Genres),
		CourseType:     course.CourseType,
		TotalSeats:     course.TotalSeats,
		SeatsBooked:    []string(course.SeatsBooked),
		AvailableSeats: course.TotalSeats - len(course.SeatsBooked),
	}
}

func (h *CourseHandler) GetAvailableCourses(c *fiber.Ctx) error {
	student := c.Locals("student").(*models.Student)

	courses, err := h.courseService.GetAvailableCourses(student.ID, student.Department)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var response []CourseResponse
	for _, course := range courses {
		response = append(response, toCourseResponse(course))
	}

	return c.JSON(fiber.Map{
		"courses": response,
	})
}

func (h *CourseHandler) GetAllCourses(c *fiber.Ctx) error {
	courses, err := h.courseService.GetAllCourses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var response []CourseResponse
	for _, course := range courses {
		response = append(response, toCourseResponse(course))
	}

	return c.JSON(fiber.Map{
		"courses": response,
	})
}

func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	var req CreateCourseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate course type
	if req.CourseType != 1 && req.CourseType != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Course type must be 1 or 2",
		})
	}

	course := &models.Course{
		Name:        req.Name,
		PDFLink:     req.PDFLink,
		Rating:      req.Rating,
		SeatsBooked: models.StringArray{},
		StaffNames:  models.StringArray(req.StaffNames),
		ImageLink:   req.ImageLink,
		Description: req.Description,
		Departments: models.StringArray(req.Departments),
		Genres:      models.StringArray(req.Genres),
		CourseType:  req.CourseType,
		TotalSeats:  req.TotalSeats,
	}

	err := h.courseService.CreateCourse(course)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Course created successfully",
		"course": fiber.Map{
			"id":          course.ID,
			"name":        course.Name,
			"course_type": course.CourseType,
			"departments": course.Departments,
			"rating":      course.Rating,
			"total_seats": course.TotalSeats,
		},
	})

}

func (h *CourseHandler) BookCourse(c *fiber.Ctx) error {
	student := c.Locals("student").(*models.Student)

	var req BookCourseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.courseService.BookCourse(student.ID, req.CourseID, req.SeatNo)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Course booked successfully",
	})
}

func (h *CourseHandler) GetMyBookings(c *fiber.Ctx) error {
	student := c.Locals("student").(*models.Student)

	bookings, err := h.courseService.GetStudentBookings(student.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"bookings": bookings,
	})
}
