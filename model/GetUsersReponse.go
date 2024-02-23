package model

type GetUsersResponse struct {
	Users map[string]User `json:"users"`
}
