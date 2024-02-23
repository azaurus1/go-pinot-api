package model

type User struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	Component             string `json:"component"`
	Role                  string `json:"role"`
	UsernameWithComponent string `json:"usernameWithComponent"`
}
