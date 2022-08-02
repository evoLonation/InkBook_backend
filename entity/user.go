package entity

type User struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Password string `json:"pwd"`
	Gender   string `json:"gender"`
	Intro    string `json:"intro"`
	Email    string `json:"email"`
}
