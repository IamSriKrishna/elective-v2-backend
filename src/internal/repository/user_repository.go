package repository

import (
	"github.com/sk/elective/src/internal/domain"
	"github.com/sk/elective/src/internal/repository/models"
	"gorm.io/gorm"
)

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) domain.StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *studentRepository) GetByRegisterNo(registerNo string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("register_no = ?", registerNo).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByID(id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}