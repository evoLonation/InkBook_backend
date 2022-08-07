package entity

import "time"

type Template struct {
	TemplateId int       `gorm:"column:template_id" json:"templateId"`
	Name       string    `gorm:"column:name" json:"name"`
	CreatorId  string    `gorm:"column:creator_id" json:"creatorId"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	Intro      string    `gorm:"column:intro" json:"intro"`
	Content    string    `gorm:"column:content" json:"content"`
}
