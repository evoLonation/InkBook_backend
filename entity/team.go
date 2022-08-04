package entity

type Team struct {
	TeamId    int    `json:"teamId"`
	Name      string `json:"teamName"`
	Intro     string `json:"teamIntroductory"`
	CaptainId string `json:"userId"`
	Url       string `json:"url"`
}
