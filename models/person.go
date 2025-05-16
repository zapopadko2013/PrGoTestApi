/* package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic"`
	Age         *int    `json:"age"`
	Gender      *string `json:"gender"`
	Nationality *string `json:"nationality"`
} */

package models

import "time"

func (Person) TableName() string {
	return "persons"
}

type Person struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Patronymic  *string   `json:"patronymic,omitempty"`
	Gender      *string   `json:"gender,omitempty"`
	Age         *int      `json:"age,omitempty"`
	Nationality *string   `json:"nationality,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PersonInput struct {
	Name       string  `json:"name" binding:"required"`
	Surname    string  `json:"surname" binding:"required"`
	Patronymic *string `json:"patronymic"`
}
