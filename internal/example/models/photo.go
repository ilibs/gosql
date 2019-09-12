package models

type Photos struct {
	Id       int    `form:"id" json:"id" db:"id"`
	MomentId int    `form:"moment_id" json:"moment_id" db:"moment_id"`
	Url      string `form:"url" json:"url" db:"url"`
	ModelTime
}

func (p *Photos) TableName() string {
	return "photos"
}

func (p *Photos) PK() string {
	return "id"
}
