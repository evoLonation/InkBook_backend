package entity

import "time"

type Document struct {
	ID         int
	Name       string
	ProjectID  int
	CreateTime time.Time
	ModifierID string
	ModifyTime time.Time
	IsEditing  bool
	IsDeleted  bool
	DeleterID  string
	DeleteTime time.Time
}
