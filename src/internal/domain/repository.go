package domain

import "github.com/sk/elective/src/internal/repository/models"

type StudentRepository interface {
	Create(student *models.Student) error
	GetByRegisterNo(registerNo string) (*models.Student, error)
	GetByID(id uint) (*models.Student, error)
}

type CourseRepository interface {
	GetAll() ([]models.Course, error)
	GetByID(id uint) (*models.Course, error)
	GetByDepartmentAndType(department string, courseType int) ([]models.Course, error)
	Update(course *models.Course) error
    Create(course *models.Course) error
}

type CourseBookingRepository interface {
	Create(booking *models.CourseBooking) error
	GetByStudentID(studentID uint) ([]models.CourseBooking, error)
	GetByStudentAndType(studentID uint, courseType int) (*models.CourseBooking, error)
	CountByStudentAndType(studentID uint, courseType int) (int64, error)
}
