package models

import (
	"github.com/ilibs/gosql"
)

type Users struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Type      int       `form:"type" json:"type" db:"type"`
	Openid    string    `form:"openid" json:"openid" db:"openid"`
	NickName  string    `form:"nickname" json:"nickname" db:"nickname"`
	Avatar    string    `form:"avatar" json:"avatar" db:"avatar"`
	City      string    `form:"city" json:"city" db:"city"`
	Country   string    `form:"country" json:"country" db:"country"`
	Gender    int       `form:"gender" json:"gender" db:"gender"`
	Province  string    `form:"province" json:"province" db:"province"`
	ModelTime
}

func (u *Users) DbName() string {
	return "default"
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) PK() string {
	return "id"
}

func GetUser(uid int) (*Users, error) {
	user := &Users{}
	err := gosql.Model(user).Where("id = ?", uid).Get()

	if err != nil {
		return nil, err
	}
	return user, nil
}

func UserGetList(start int, num int) ([]*Users, error) {
	var m = make([]*Users, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}
