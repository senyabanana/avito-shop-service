package entity

type User struct {
	ID       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	Coins    int    `json:"coins"`
}
