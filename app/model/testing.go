package model

import "testing"

//TestUser ...
func TestUser(t *testing.T) *User {
	return &User{
	ID: 55555,
	Login: "Lev",
	Username: "Lev",
	Surname: "Lvovich",
	Birthday: "13.09.1994",
	Password: "LevIsTheBest",
}
}