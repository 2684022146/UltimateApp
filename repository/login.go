package repository

import (
	"context"
	"errors"
	"fmt"
	"webdemo/model"
	"webdemo/util"

	"gorm.io/gorm"
)

type LoginRepository interface {
	Login(ctx context.Context, username, password string, roleId int8) (*model.User, error)
	Regist(ctx context.Context, username, password string, phone string, roleId int8) error
}
type loginRepository struct {
	db gorm.DB
}

func NewLoginRepository(db *gorm.DB) LoginRepository {
	return &loginRepository{
		db: *db,
	}
}
func (r *loginRepository) Login(ctx context.Context, username, password string, roleId int8) (*model.User, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("用户名和密码不能为空")
	}

	hashedPassword := util.Md5String(password)
	var user *model.User
	err := r.db.WithContext(ctx).Table("users").Where("username=? AND password=? AND role_id=?", username, hashedPassword, roleId).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("username or password error")
		}
		return nil, fmt.Errorf("query user fial")
	}
	return user, nil
}
func (r *loginRepository) Regist(ctx context.Context, username, password string, phone string, roleId int8) error {
	if username == "" || password == "" {
		return fmt.Errorf("用户名和密码不能为空")
	}
	if len(password) < 6 || len(password) > 8 {
		return fmt.Errorf("密码长度在6位到8位间")
	}
	if len(phone) != 11 {
		return fmt.Errorf("手机号错误")
	}

	var isDuplicate bool
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username=? AND role_id=?", username, roleId).Select("1").Limit(1).Scan(&isDuplicate).Error
	if err != nil {
		return fmt.Errorf("check duplicate username fail:%v", err)
	}
	if isDuplicate {
		return fmt.Errorf("username is dulicate")
	}
	hashedPassword := util.Md5String(password)
	newUser := &model.User{
		Username: username,
		Password: hashedPassword,
		Phone:    phone,
		RoleID:   roleId,
	}
	var txErr error
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if txErr == nil {
			if commitErr := tx.Commit(); commitErr != nil {
				txErr = fmt.Errorf("提交事务失败")
				return
			}
		}
	}()
	err = tx.WithContext(ctx).Model(&model.User{}).Create(newUser).Error
	if err != nil {
		txErr = fmt.Errorf("create new user fail:%v", err)
		return txErr
	}
	return txErr
}
