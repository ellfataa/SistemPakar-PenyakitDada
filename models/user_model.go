package models

type User struct {
	IDUser   int    `json:"id_user"`
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

type RegisterRequest struct {
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRequest struct {
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}