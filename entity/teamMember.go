package entity

type Member struct {
	TeamId   int    `json:"teamId"`
	MemberId string `json:"memberId"`
	Identity int    `json:"identity"`
}
