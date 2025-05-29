package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, a)
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}
	return json.Marshal(a)
}

type Student struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	RegisterNo string    `json:"register_no" gorm:"unique;not null"`
	Password   string    `json:"-" gorm:"not null"`
	Name       string    `json:"name" gorm:"not null"`
	Department string    `json:"department" gorm:"default:'CSE'"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Bookings
	CourseBookings []CourseBooking `json:"course_bookings,omitempty" gorm:"foreignKey:StudentID"`
}
type Course struct {
	ID             uint        `json:"id" gorm:"primaryKey"`
	Name           string      `json:"name" gorm:"not null"`
	PDFLink        string      `json:"pdf_link"`
	Rating         float64     `json:"rating" gorm:"default:0"`
	SeatsBooked    StringArray `json:"seats_booked" gorm:"type:jsonb"`
	StaffNames     StringArray `json:"staff_names" gorm:"type:jsonb"`
	ImageLink      string      `json:"image_link"`
	Description    string      `json:"description"`
	Departments    StringArray `json:"departments" gorm:"type:jsonb"`
	Genres         StringArray `json:"genres" gorm:"type:jsonb"`
	CourseType     int         `json:"course_type" gorm:"not null"`
	TotalSeats     int         `json:"total_seats" gorm:"not null"`
	AvailableSeats int         `json:"available_seats"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`

	CourseBookings []CourseBooking `json:"course_bookings,omitempty" gorm:"foreignKey:CourseID"`
}

type CourseBooking struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StudentID uint      `json:"student_id" gorm:"not null"`
	CourseID  uint      `json:"course_id" gorm:"not null"`
	SeatNo    string    `json:"seat_no"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Student Student `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	Course  Course  `json:"course,omitempty" gorm:"foreignKey:CourseID"`
}

type StudentEntity struct {
	ID         uint   `json:"id"`
	RegisterNo string `json:"register_no"`
	Department string `json:"department"`
	Name       string `json:"name"`
}

type CourseEntity struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	PDFLink     string   `json:"pdf_link"`
	Rating      float64  `json:"rating"`
	SeatsBooked []string `json:"seats_booked"`
	StaffNames  []string `json:"staff_names"`
	ImageLink   string   `json:"image_link"`
	Description string   `json:"description"`
	Departments []string `json:"departments"`
	Genres      []string `json:"genres"`
	CourseType  int      `json:"course_type"`
}
