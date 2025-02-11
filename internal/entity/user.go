package entity

type User struct {
	ID       int    `json:"-" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password_hash"`
	Coins    int    `json:"coins" db:"coins"`
}
