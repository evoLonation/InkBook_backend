package entity

import "time"

type Folder struct {
	FolderId   int       `gorm:"column:folder_id" json:"folderId"`
	Name       string    `gorm:"column:name" json:"name"`
	TeamId     int       `gorm:"column:team_id" json:"teamId"`
	ParentId   int       `gorm:"column:parent_id" json:"parentId"`
	CreatorId  string    `gorm:"column:creator_id" json:"creatorId"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	IsDeleted  bool      `gorm:"column:is_deleted" json:"isDeleted"`
	DeleterId  string    `gorm:"column:deleter_id" json:"deleterId"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"deleteTime"`
}
