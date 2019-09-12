package models

type Moments struct {
	Id           int    `form:"id" json:"id" db:"id"`
	UserId       int    `form:"user_id" json:"user_id" db:"user_id"`
	Content      string `form:"content" json:"content" db:"content"`
	CommentTotal int    `form:"comment_total" json:"comment_total" db:"comment_total"`
	LikeTotal    int    `form:"like_total" json:"like_total" db:"like_total"`
	Status       int    `form:"status" json:"status" db:"status"`
	ModelTime
}

func (p *Moments) TableName() string {
	return "moments"
}

func (p *Moments) PK() string {
	return "id"
}
