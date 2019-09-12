package models

import (
	"database/sql"
	"time"
)

type Users struct {
	Id          int            `form:"id" json:"id" db:"id"`
	Name        string         `form:"name" json:"name" db:"name"`
	Status      int            `form:"status" json:"status" db:"status"`
	SuccessTime sql.NullString `form:"-" json:"success_time" db:"success_time"`
	CreatedAt   time.Time      `form:"-" json:"created_at" db:"created_at" time_format:"2006-01-02 15:04:05"`
	UpdatedAt   time.Time      `form:"-" json:"updated_at" db:"updated_at" time_format:"2006-01-02 15:04:05"`
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) PK() string {
	return "id"
}
