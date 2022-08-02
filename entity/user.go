package entity

type User struct {
	UserId   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Intro    string `json:"intro"`
	Email    string `json:"email"`
}
