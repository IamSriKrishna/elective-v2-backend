package domain

import "github.com/sk/elective/src/internal/repository/models"

type AuthService interface {
    Register(registerNo, password, department string) (*models.Student, error)
    Login(registerNo, password string) (string, *models.Student, error)
    ValidateToken(token string) (*models.Student, error)
}

type CourseService interface {
    GetAvailableCourses(studentID uint, department string) ([]models.Course, error)
    BookCourse(studentID uint, courseID uint, seatNo string) error
    GetStudentBookings(studentID uint) ([]models.CourseBooking, error)
    CreateCourse(course *models.Course) error 
    GetAllCourses()  ([]models.Course, error)
}