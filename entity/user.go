package entity

type User struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Nickname string `gorm:"size:20;not null"`
	RealName string `gorm:"size:20;not null"`
	Password string `gorm:"size:20;not null"`
	Gender   string
	Intro    string `gorm:"size:255"`
	Email    string `gorm:"unique_index;not null"`
}
