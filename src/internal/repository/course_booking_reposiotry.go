package repository

import (
	"github.com/sk/elective/src/internal/domain"
	"github.com/sk/elective/src/internal/repository/models"
	"gorm.io/gorm"
)

type courseBookingRepository struct {
    db *gorm.DB
}

func NewCourseBookingRepository(db *gorm.DB) domain.CourseBookingRepository {
    return &courseBookingRepository{db: db}
}

func (r *courseBookingRepository) Create(booking *models.CourseBooking) error {
    return r.db.Create(booking).Error
}

func (r *courseBookingRepository) GetByStudentID(studentID uint) ([]models.CourseBooking, error) {
    var bookings []models.CourseBooking
    err := r.db.Preload("Course").Where("student_id = ?", studentID).Find(&bookings).Error
    return bookings, err
}

func (r *courseBookingRepository) GetByStudentAndType(studentID uint, courseType int) (*models.CourseBooking, error) {
    var booking models.CourseBooking
    err := r.db.Joins("JOIN courses ON courses.id = course_bookings.course_id").
        Where("course_bookings.student_id = ? AND courses.course_type = ?", studentID, courseType).
        First(&booking).Error
    if err != nil {
        return nil, err
    }
    return &booking, nil
}

func (r *courseBookingRepository) CountByStudentAndType(studentID uint, courseType int) (int64, error) {
    var count int64
    err := r.db.Table("course_bookings").
        Joins("JOIN courses ON courses.id = course_bookings.course_id").
        Where("course_bookings.student_id = ? AND courses.course_type = ?", studentID, courseType).
        Count(&count).Error
    return count, err
}
