package models

import (
	"github.com/ilibs/gosql"
)

type Photos struct {
	Id        int       `form:"id" json:"id" db:"id"`
	MomentId  int       `form:"moment_id" json:"moment_id" db:"moment_id"`
	Url       string    `form:"url" json:"url" db:"url"`
	ModelTime
}

func (p *Photos) DbName() string {
	return "default"
}

func (p *Photos) TableName() string {
	return "photos"
}

func (p *Photos) PK() string {
	return "id"
}

func GetPhotos(id int) ([]*Photos, error) {
	var m = make([]*Photos, 0)
	err := gosql.Model(&m).Where("moment_id = ?", id).OrderBy("id desc").All()
	if err != nil {
		return nil, err
	}
	return m, nil
}
