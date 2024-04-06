package user

import "context"

type User struct {
	ID        int64  `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Firstname string `json:"firstname" db:"first_name"`
	Lastname  string `json:"lastname" db:"last_name"`
}

type CreateUserReq struct {
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Firstname string `json:"firstname" db:"first_name"`
	Lastname  string `json:"lastname" db:"last_name"`
}

type CreateUserRes struct {
	ID        string `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Firstname string `json:"firstname" db:"first_name"`
	Lastname  string `json:"lastname" db:"last_name"`
}

type LoginUserReq struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type LoginUserRes struct {
	accessToken string
	ID          string `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	Firstname   string `json:"firstname" db:"first_name"`
	Lastname    string `json:"lastname" db:"last_name"`
}

type Repo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
}
