package entity

import (
	"time"
)

type Graph struct {
	GraphId    int       `gorm:"column:graph_id" json:"graphId"`
	Name       string    `gorm:"column:name" json:"name"`
	ProjectId  int       `gorm:"column:project_id" json:"projectId"`
	CreatorId  string    `gorm:"column:creator_id" json:"creatorId"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	ModifierId string    `gorm:"column:modifier_id" json:"modifierId"`
	ModifyTime time.Time `gorm:"column:modify_time" json:"modifyTime"`
	IsEditing  bool      `gorm:"column:is_editing" json:"isEditing"`
	IsDeleted  bool      `gorm:"column:is_deleted" json:"isDeleted"`
	DeleterId  string    `gorm:"column:deleter_id" json:"deleterId"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"deleteTime"`
	Content    string    `gorm:"column:content" json:"content"`
	EditingCnt int       `gorm:"column:editing_cnt" json:"editingCnt"`
	Img        string    `gorm:"column:img" json:"img"`
}
