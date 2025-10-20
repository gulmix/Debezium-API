package models

type UsersDTO struct {
	Users []UserDTO `json:"users"`
	Total int       `json:"total"`
}

type UserDTO struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}
