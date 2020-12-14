package service

type UserRequest struct {
	Uid int `json:"uid"`
	Method string
}

type UserResponse struct {
	Name string `json:"name"`
}

