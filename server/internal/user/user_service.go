package user

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/golang-jwt/jwt/v4"
	"server/util"
	"strconv"
	"time"
)

type service struct {
	Repo
	timeout time.Duration
}

func NewService(repository Repo) Service {
	return &service{
		repository,
		time.Duration(5) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Password:  hashedPassword,
	}

	r, err := s.Repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		ID:        strconv.Itoa(int(r.ID)),
		Username:  r.Username,
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
	}

	return res, nil
}

type CustomClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return &LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	ss, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		return &LoginUserRes{}, err
	}

	return &LoginUserRes{accessToken: ss, ID: strconv.Itoa(int(u.ID)), Username: u.Username, Firstname: u.Firstname, Lastname: u.Lastname}, nil
}
