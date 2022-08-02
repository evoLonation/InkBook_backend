package entity

type Team struct {
	ID        int    `json:"teamId"`
	Name      string `json:"teamName"`
	Intro     string `json:"teamIntroductory"`
	CaptainID string `json:"userId"`
	Url       string `json:"url"`
}
