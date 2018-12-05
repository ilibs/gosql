package models

import (
	"github.com/ilibs/gosql"
)

type Moments struct {
	Id           int       `form:"id" json:"id" db:"id"`
	UserId       int       `form:"user_id" json:"user_id" db:"user_id"`
	Content      string    `form:"content" json:"content" db:"content"`
	CommentTotal int       `form:"comment_total" json:"comment_total" db:"comment_total"`
	LikeTotal    int       `form:"like_total" json:"like_total" db:"like_total"`
	Status       int       `form:"status" json:"status" db:"status"`
	ModelTime
}

func (p *Moments) DbName() string {
	return "default"
}

func (p *Moments) TableName() string {
	return "moments"
}

func (p *Moments) PK() string {
	return "id"
}

type UserMoment struct {
	Moments
	User   *Users    `json:"user" db:"-" relation:"user_id,id"`
	Photos []*Photos `json:"photos" db:"-" relation:"id,moment_id"`
}

func MomentGet(id int) (*UserMoment, error) {
	moment := &UserMoment{}
	err := gosql.Model(moment).Where("status = 1 and id = ?",id).Get()

	if err != nil {
		return nil, err
	}

	return moment, err
}

func MomentGetList() ([]*UserMoment, error) {
	var moments = make([]*UserMoment, 0)
	err := gosql.Model(&moments).Where("status = 1").Limit(10).All()
	if err != nil {
		return nil, err
	}
	return moments, err
}
