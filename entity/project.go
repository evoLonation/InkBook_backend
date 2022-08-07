package entity

import "time"

type Project struct {
	ProjectId  int       `gorm:"column:project_id" json:"projectId"`
	TeamId     int       `gorm:"column:team_id" json:"teamId"`
	Name       string    `gorm:"column:name" json:"name"`
	CreatorId  string    `gorm:"column:creator_id" json:"creatorId"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	IsDeleted  bool      `gorm:"column:is_deleted" json:"isDeleted"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"deleteTime"`
	Intro      string    `gorm:"column:intro" json:"intro"`
	ImgURL     string    `gorm:"column:img_url" json:"imgUrl"`
}
