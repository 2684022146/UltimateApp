package service

import (
	"context"
	"fmt"
	"webdemo/model"
	"webdemo/repository"
)

type SettingsService interface {
	CreateAddress(ctx context.Context, req *model.CreateAddressRequest, userId uint) error
	AddressList(ctx context.Context, userId uint) ([]*model.Address, error)
	AddressDetail(ctx context.Context, addressId, userId uint) (*model.Address, error)
	UpdateAddress(ctx context.Context, req *model.Address) error
	DeleteAddress(ctx context.Context, addressId, userId uint) error
	SetDefault(ctx context.Context, addressId, userId uint) error
}
type settingsService struct {
	repo repository.SettingsRepository
}

func NewSettingsService(repo repository.SettingsRepository) SettingsService {
	return &settingsService{
		repo: repo,
	}
}
func (s *settingsService) CreateAddress(ctx context.Context, req *model.CreateAddressRequest, userId uint) error {
	address := &model.Address{
		UserID:    userId,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Street:    req.Street,
		Detail:    req.Detail,
		Receiver:  req.Receiver,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		IsDefault: req.IsDefault,
	}
	err := s.repo.CreateAddress(ctx, address)
	if err != nil {
		return fmt.Errorf("创建新地址失败:%w", err)
	}
	return nil
}
func (s *settingsService) AddressList(ctx context.Context, userId uint) ([]*model.Address, error) {
	addressSlice, err := s.repo.AddressList(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return addressSlice, nil
}
func (s *settingsService) AddressDetail(ctx context.Context, addressId, userId uint) (*model.Address, error) {
	return s.repo.AddressDetail(ctx, addressId, userId)
}
func (s *settingsService) UpdateAddress(ctx context.Context, req *model.Address) error {
	if err := s.repo.UpdateAddress(ctx, req); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
func (s *settingsService) DeleteAddress(ctx context.Context, addressId, userId uint) error {
	err := s.repo.DeleteAddress(ctx, addressId, userId)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return err
}
func (s *settingsService) SetDefault(ctx context.Context, addressId, userId uint) error {
	err := s.repo.SetDefault(ctx, addressId, userId)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}
