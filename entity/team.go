package entity

type Team struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name      string `gorm:"size:20;not null"`
	Intro     string `gorm:"size:255"`
	CaptainID uint
	Captain   User `gorm:"ForeignKey:ID;references:CaptainID"`
}
