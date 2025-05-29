package repository

import (
	"github.com/sk/elective/src/internal/domain"
	"github.com/sk/elective/src/internal/repository/models"
	"gorm.io/gorm"
)

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) domain.CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Find(&courses).Error
	return courses, err
}

func (r *courseRepository) GetByID(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) GetByDepartmentAndType(department string, courseType int) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Where("course_type = ? AND departments @> ?", courseType, `["`+department+`"]`).Find(&courses).Error
	return courses, err
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Create(course *models.Course) error {
    return r.db.Create(course).Error
}