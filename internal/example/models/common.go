package models

import (
	"time"
)

type ModelTime struct {
	CreatedAt time.Time `form:"-" json:"created_at" db:"created_at" time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `form:"-" json:"updated_at" db:"updated_at" time_format:"2006-01-02 15:04:05"`
}
