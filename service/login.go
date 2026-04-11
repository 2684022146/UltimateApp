package service

import (
	"context"
	"fmt"
	"time"
	"webdemo/model"
	"webdemo/repository"
	"webdemo/util"
)

type LoginService interface { //接口
	Login(ctx context.Context, req *model.LoginRequest) (string, error)
	Regist(ctx context.Context, req *model.RegisterRequest) (string, error)
}
type loginService struct {
	repo repository.LoginRepository
}

func NewLoginService(repo repository.LoginRepository) LoginService {
	return &loginService{
		repo: repo,
	}
}
func (s *loginService) Login(ctx context.Context, req *model.LoginRequest) (string, error) {
	if req.Username == "" || req.Password == "" {
		return "", fmt.Errorf("username and password not empty")
	}
	user, err := s.repo.Login(ctx, req.Username, req.Password, req.RoleID)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	token, err := util.GenerateToken(user.Id, user.Username, user.RoleID, time.Hour*2)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	return token, nil
}
func (s *loginService) Regist(ctx context.Context, req *model.RegisterRequest) (string, error) {
	if req.Username == "" || req.Password == "" {
		return "", fmt.Errorf("username and password not empty")
	}
	if req.Phone == "" {
		return "", fmt.Errorf("phone not empty")
	}
	if err := s.repo.Regist(ctx, req.Username, req.Password, req.Phone, req.RoleID); err != nil {
		return "", fmt.Errorf("%v", err)
	}
	loginRequest := &model.LoginRequest{
		Username: req.Username,
		Password: req.Password,
		RoleID:   req.RoleID,
	}
	token, err := s.Login(ctx, loginRequest)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	return token, nil
}
