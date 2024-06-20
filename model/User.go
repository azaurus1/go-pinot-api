package model

type User struct {
	Username              string    `json:"username"`
	Password              string    `json:"password"`
	Component             string    `json:"component"`
	Role                  string    `json:"role"`
	UsernameWithComponent string    `json:"usernameWithComponent"`
	Permissions           *[]string `json:"permissions,omitempty"` // optional for adding permissions etc
	Tables                *[]string `json:"tables,omitempty"`
}
