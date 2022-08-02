package entity

import "time"

type Project struct {
	ProjectID  int       `json:"projectId"`
	TeamID     int       `json:"teamId"`
	Name       string    `json:"name"`
	CreatorID  int       `json:"creatorId"`
	CreateTime time.Time `json:"createTime"`
	IsDeleted  bool      `json:"isDeleted"`
	DeleteTime time.Time `json:"deleteTime"`
	Intro      string    `json:"intro"`
}
