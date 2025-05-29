package usecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sk/elective/src/internal/config"
	"github.com/sk/elective/src/internal/domain"
	"github.com/sk/elective/src/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	studentRepo domain.StudentRepository
	jwtConfig   config.JWTConfig
}

func NewAuthService(studentRepo domain.StudentRepository, jwtConfig config.JWTConfig) domain.AuthService {
	return &authService{
		studentRepo: studentRepo,
		jwtConfig:   jwtConfig,
	}
}

type Claims struct {
	StudentID  uint   `json:"student_id"`
	RegisterNo string `json:"register_no"`
	Department string `json:"department"`
	Name       string `json:"name"`
	jwt.RegisteredClaims
}

func (s *authService) Register(registerNo, password, department, name string) (*models.Student, error) {
	// Check if student already exists
	existingStudent, err := s.studentRepo.GetByRegisterNo(registerNo)
	if err == nil && existingStudent != nil {
		return nil, errors.New("student already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create student
	student := &models.Student{
		RegisterNo: registerNo,
		Password:   string(hashedPassword),
		Department: department,
		Name:       name,
	}

	err = s.studentRepo.Create(student)
	if err != nil {
		return nil, err
	}

	return student, nil
}

func (s *authService) Login(registerNo, password string) (string, *models.Student, error) {
	student, err := s.studentRepo.GetByRegisterNo(registerNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	claims := &Claims{
		StudentID:  student.ID,
		RegisterNo: student.RegisterNo,
		Department: student.Department,
		Name:       student.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtConfig.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, student, nil
}

func (s *authService) ValidateToken(tokenString string) (*models.Student, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	student, err := s.studentRepo.GetByID(claims.StudentID)
	if err != nil {
		return nil, err
	}

	return student, nil
}
