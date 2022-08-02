package entity

import "time"

type Project struct {
	ID         int
	TeamID     int
	Name       string
	CreateTime time.Time
	IsDeleted  bool
	DeleteTime time.Time
}
