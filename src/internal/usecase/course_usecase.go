package usecase

import (
	"errors"
	"fmt"

	"github.com/sk/elective/src/internal/domain"
	"github.com/sk/elective/src/internal/repository/models"
)

type courseService struct {
	courseRepo  domain.CourseRepository
	bookingRepo domain.CourseBookingRepository
}

func NewCourseService(courseRepo domain.CourseRepository, bookingRepo domain.CourseBookingRepository) domain.CourseService {
	return &courseService{
		courseRepo:  courseRepo,
		bookingRepo: bookingRepo,
	}
}
func (s *courseService) GetAllCourses() ([]models.Course, error) {
	return s.courseRepo.GetAll()
}
func (s *courseService) GetAvailableCourses(studentID uint, department string) ([]models.Course, error) {
	var availableCourses []models.Course

	// Get type 1 courses
	type1Courses, err := s.courseRepo.GetByDepartmentAndType(department, 1)
	if err != nil {
		return nil, err
	}

	// Check if student has already booked a type 1 course
	type1Count, err := s.bookingRepo.CountByStudentAndType(studentID, 1)
	if err != nil {
		return nil, err
	}

	if type1Count == 0 {
		availableCourses = append(availableCourses, type1Courses...)
	}

	// Get type 2 courses
	type2Courses, err := s.courseRepo.GetByDepartmentAndType(department, 2)
	if err != nil {
		return nil, err
	}

	// Check if student has already booked a type 2 course
	type2Count, err := s.bookingRepo.CountByStudentAndType(studentID, 2)
	if err != nil {
		return nil, err
	}

	if type2Count == 0 {
		availableCourses = append(availableCourses, type2Courses...)
	}

	return availableCourses, nil
}

func (s *courseService) BookCourse(studentID uint, courseID uint, seatNo string) error {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return err
	}

	// Check if student has already booked a course of this type
	count, err := s.bookingRepo.CountByStudentAndType(studentID, course.CourseType)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(fmt.Sprintf("you have already booked a type %d course", course.CourseType))
	}

	// Check if seat is already booked
	for _, bookedSeat := range course.SeatsBooked {
		if bookedSeat == seatNo {
			return errors.New("seat already booked")
		}
	}

	// Create booking
	booking := &models.CourseBooking{
		StudentID: studentID,
		CourseID:  courseID,
		SeatNo:    seatNo,
	}

	err = s.bookingRepo.Create(booking)
	if err != nil {
		return err
	}

	// Update course seats
	course.SeatsBooked = append(course.SeatsBooked, seatNo)
	err = s.courseRepo.Update(course)
	if err != nil {
		return err
	}

	return nil
}

func (s *courseService) GetStudentBookings(studentID uint) ([]models.CourseBooking, error) {
	return s.bookingRepo.GetByStudentID(studentID)
}

func (s *courseService) CreateCourse(course *models.Course) error {
	// Validate course type
	if course.CourseType != 1 && course.CourseType != 2 {
		return errors.New("course type must be 1 or 2")
	}

	// Validate required fields
	if course.Name == "" {
		return errors.New("course name is required")
	}

	// Initialize empty arrays if nil
	if course.SeatsBooked == nil {
		course.SeatsBooked = models.StringArray{}
	}
	if course.StaffNames == nil {
		course.StaffNames = models.StringArray{}
	}
	if course.Departments == nil {
		course.Departments = models.StringArray{}
	}
	if course.Genres == nil {
		course.Genres = models.StringArray{}
	}

	return s.courseRepo.Create(course)
}
