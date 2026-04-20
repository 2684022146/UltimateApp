package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

type GeocodeResponse struct {
	Status   string    `json:"status"`
	Geocodes []Geocode `json:"geocodes"`
}
type Geocode struct {
	Location string `json:"location"`
}

const (
	apiUrl = "https://restapi.amap.com/v3/geocode/geo?address=%s&output=JSON&key=74e8efb95dcff3e269ce54497a8d1b18"
)

func (s *settingsService) CreateAddress(ctx context.Context, req *model.CreateAddressRequest, userId uint) error {
	completeAddress := fmt.Sprintf("%s%s%s%s%s", req.Province, req.City, req.District, req.Street, req.Detail)
	url := fmt.Sprintf(apiUrl, completeAddress)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("api失败:%w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取地址信息失败:%w", err)
	}
	fmt.Println(string(body))

	var geocodeResp GeocodeResponse
	err = json.Unmarshal(body, &geocodeResp)
	if err != nil {
		return fmt.Errorf("解析地址信息失败:%w", err)
	}
	if geocodeResp.Status != "1" || len(geocodeResp.Geocodes) == 0 {
		return fmt.Errorf("api获取geo失败:%s", geocodeResp.Status)
	}
	coordinates := geocodeResp.Geocodes[0].Location
	parts := strings.Split(coordinates, ",")
	longitudeStr := parts[0]
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitudeStr := parts[1]
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return fmt.Errorf("获取地址信息失败:%w", err)
	}
	defer resp.Body.Close()
	address := &model.Address{
		UserID:    userId,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Street:    req.Street,
		Detail:    req.Detail,
		Receiver:  req.Receiver,
		Phone:     req.Phone,
		Latitude:  latitude,
		Longitude: longitude,
		IsDefault: req.IsDefault,
	}
	err = s.repo.CreateAddress(ctx, address)
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
