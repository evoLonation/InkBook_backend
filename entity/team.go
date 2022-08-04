package entity

type Team struct {
	TeamId    int    `gorm:"team_id" json:"teamId"`
	Name      string `json:"teamName"`
	Intro     string `json:"teamIntroductory"`
	CaptainId string `json:"userId"`
	Url       string `json:"url"`
}
