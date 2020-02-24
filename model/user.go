package model

type User struct {
	UserID    int    `db:"user_id"`
	UserName  string `db:"user_name"`
	Email     string `db:"email"`
	Telephone string `db:"telephone"`
}
