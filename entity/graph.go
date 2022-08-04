package entity

import (
	"time"
)

type Graph struct {
	GraphID    int       `gorm:"column:graph_id" json:"graphId"`
	Name       string    `gorm:"column:name" json:"name"`
	ProjectID  int       `gorm:"column:project_id" json:"projectId"`
	CreatorID  string    `gorm:"column:creator_id" json:"creatorId"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	ModifierID string    `gorm:"column:modifier_id" json:"modifierId"`
	ModifyTime time.Time `gorm:"column:modify_time" json:"modifyTime"`
	IsEditing  bool      `gorm:"column:is_editing" json:"isEditing"`
	IsDeleted  bool      `gorm:"column:is_deleted" json:"isDeleted"`
	DeleterID  string    `gorm:"column:deleter_id" json:"deleterId"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"deleteTime"`
	Content    string    `gorm:"column:content" json:"content"`
}
