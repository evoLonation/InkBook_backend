package entity

import "time"

type Project struct {
	ProjectID  int       `gorm:"column:project_id" json:"projectId"`
	TeamID     int       `gorm:"column:team_id" json:"teamId"`
	Name       string    `gorm:"column:name" json:"name"`
	CreatorID  string    `gorm:"column:creator_id" json:"creatorId"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	IsDeleted  bool      `gorm:"column:is_deleted" json:"isDeleted"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"deleteTime"`
	Intro      string    `gorm:"column:intro" json:"intro"`
	ImgURL     string    `gorm:"column:img_url" json:"imgUrl"`
}
