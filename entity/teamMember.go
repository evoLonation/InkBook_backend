package entity

type TeamMember struct {
	TeamId   int    `json:"teamId"`
	MemberId string `json:"memberId"`
	Identity int    `json:"identity"`
}
