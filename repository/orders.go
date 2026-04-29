package repository

import (
	"context"
	"log"
	"time"
	"webdemo/consts"
	"webdemo/model"

	"gorm.io/gorm"
)

type OrdersRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	OrderDetailBasic(ctx context.Context, orderNo string) (*model.Order, error)
	//
	SenderFinishedOrder(ctx context.Context, userID uint, pageInfo *model.Page) ([]*model.Order, int64, error)
	SenderInTransitOrder(ctx context.Context, userID uint, pageInfo *model.Page) ([]*model.Order, int64, error)
	SenderWaitingOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.Order, int64, error)
	SenderCencelOrder(ctx context.Context, orderNo string) error
	//
	ReceiverFinishedOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.Order, int64, error)
	ReceiverInTransitOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.Order, int64, error)
	//
	RiderOrderList(ctx context.Context, status int, pageInfo *model.Page) ([]*model.Order, int64, error)
	// 骑手相关
	AcceptOrder(ctx context.Context, orderNo string, deliveryUserID uint) error
	PickupOrder(ctx context.Context, orderNo string) error
	StartDelivery(ctx context.Context, orderNo string) error
	UploadLocation(ctx context.Context, orderID uint, deliveryUserID uint, longitude, latitude float64) error
	CompleteOrder(ctx context.Context, orderNo string) error
}
type ordersRepository struct {
	db *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) OrdersRepository {
	return &ordersRepository{
		db: db,
	}
}
func (r *ordersRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(order).Error
		if err != nil {
			return err
		}
		return nil
	})
}
func (r *ordersRepository) OrderDetailBasic(ctx context.Context, orderNo string) (*model.Order, error) {
	var order model.Order

	// 1. 先查询订单基本信息
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}

	// 2. 再根据sender_address_id查询地址信息
	var address model.Address
	err = r.db.WithContext(ctx).Where("id = ?", order.SenderAddressID).First(&address).Error
	if err != nil {
		return nil, err
	}

	// 3. 组装寄件人地址信息到order结构体
	order.SenderProvince = address.Province
	order.SenderCity = address.City
	order.SenderDistrict = address.District
	order.SenderStreet = address.Street
	order.SenderDetail = address.Detail
	order.SenderReceiver = address.Receiver
	order.SenderPhone = address.Phone
	order.SenderLatitude = address.Latitude
	order.SenderLongitude = address.Longitude

	return &order, nil
}

// ///////////////////
func (r *ordersRepository) SenderFinishedOrder(ctx context.Context, userID uint, pageInfo *model.Page) ([]*model.Order, int64, error) {
	var orderList []*model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).
		Joins("LEFT JOIN addresses ON orders.sender_address_id = addresses.id").
		Where("orders.sender_user_id=? and orders.status=?", userID, consts.OrderFinished).
		Select("orders.*, addresses.province as sender_province, addresses.city as sender_city, addresses.district as sender_district, addresses.street as sender_street, addresses.detail as sender_detail, addresses.receiver as sender_receiver, addresses.phone as sender_phone, addresses.latitude as sender_latitude, addresses.longitude as sender_longitude")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	currentPage := pageInfo.CurrentPage
	if currentPage < 1 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageInfo.PerPage
	log.Printf("offset: %d, limit: %d", offset, pageInfo.PerPage)
	err := query.Offset(offset).Limit(pageInfo.PerPage).Find(&orderList).Error
	if err != nil {
		return nil, total, err
	}
	return orderList, total, nil
}

func (r *ordersRepository) SenderInTransitOrder(ctx context.Context, userID uint, pageInfo *model.Page) ([]*model.Order, int64, error) {
	var orderList []*model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).
		Joins("LEFT JOIN addresses ON orders.sender_address_id = addresses.id").
		Where("orders.sender_user_id=? and orders.status IN (?)", userID, []int{consts.OrderBeforePickUp, consts.OrderPickup, consts.OrderDelivering}).
		Select("orders.*, addresses.province as sender_province, addresses.city as sender_city, addresses.district as sender_district, addresses.street as sender_street, addresses.detail as sender_detail, addresses.receiver as sender_receiver, addresses.phone as sender_phone, addresses.latitude as sender_latitude, addresses.longitude as sender_longitude")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	currentPage := pageInfo.CurrentPage
	if currentPage < 1 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageInfo.PerPage
	err := query.Offset(offset).Limit(pageInfo.PerPage).Find(&orderList).Error
	if err != nil {
		return nil, total, err
	}
	return orderList, total, nil
}
func (r *ordersRepository) SenderWaitingOrder(ctx context.Context, usrId uint, pageInfo *model.Page) ([]*model.Order, int64, error) {
	var orderList []*model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).
		Joins("LEFT JOIN addresses ON orders.sender_address_id = addresses.id").
		Where("orders.sender_user_id=? and orders.status=?", usrId, consts.OrderWaiting).
		Select("orders.*, addresses.province as sender_province, addresses.city as sender_city, addresses.district as sender_district, addresses.street as sender_street, addresses.detail as sender_detail, addresses.receiver as sender_receiver, addresses.phone as sender_phone, addresses.latitude as sender_latitude, addresses.longitude as sender_longitude")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	log.Printf("pageInfo.per_page: %d", pageInfo.PerPage)

	currentPage := pageInfo.CurrentPage
	if currentPage < 1 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageInfo.PerPage
	err := query.Offset(offset).Limit(pageInfo.PerPage).Find(&orderList).Error
	if err != nil {
		return nil, total, err
	}
	return orderList, total, nil
}
func (r *ordersRepository) SenderCencelOrder(ctx context.Context, orderNo string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Order{}).Where("order_no=?", orderNo).Update("status", consts.OrderCancel).Error
	})
	return err
}

// ////////////////////////
func (r *ordersRepository) ReceiverFinishedOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.Order, int64, error) {
	var orderList []*model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).
		Joins("JOIN users ON orders.receiver_phone=users.phone").
		Joins("LEFT JOIN addresses ON orders.sender_address_id = addresses.id").
		Where("users.id=? and orders.status=?", userId, consts.OrderFinished).
		Select("orders.*, addresses.province as sender_province, addresses.city as sender_city, addresses.district as sender_district, addresses.street as sender_street, addresses.detail as sender_detail, addresses.receiver as sender_receiver, addresses.phone as sender_phone, addresses.latitude as sender_latitude, addresses.longitude as sender_longitude")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	currentPage := pageInfo.CurrentPage
	if currentPage < 1 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageInfo.PerPage
	err := query.Offset(offset).Limit(pageInfo.PerPage).Find(&orderList).Error
	if err != nil {
		return nil, total, err
	}
	return orderList, total, nil
}
func (r *ordersRepository) ReceiverInTransitOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.Order, int64, error) {
	var orderList []*model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).
		Joins("JOIN users ON orders.receiver_phone=users.phone").
		Joins("LEFT JOIN addresses ON orders.sender_address_id = addresses.id").
		Where("users.id=? and orders.status=?", userId, []int{consts.OrderBeforePickUp, consts.OrderPickup, consts.OrderDelivering}).
		Select("orders.*, addresses.province as sender_province, addresses.city as sender_city, addresses.district as sender_district, addresses.street as sender_street, addresses.detail as sender_detail, addresses.receiver as sender_receiver, addresses.phone as sender_phone, addresses.latitude as sender_latitude, addresses.longitude as sender_longitude")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	currentPage := pageInfo.CurrentPage
	if currentPage < 1 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageInfo.PerPage
	err := query.Offset(offset).Limit(pageInfo.PerPage).Find(&orderList).Error
	if err != nil {
		return nil, total, err
	}
	return orderList, total, nil
}

// ///////////////////
func (r *ordersRepository) RiderOrderList(ctx context.Context, status int, pageInfo *model.Page) ([]*model.Order, int64, error) {
	var orderList []*model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).
		Joins("LEFT JOIN addresses ON orders.sender_address_id = addresses.id").
		Where("orders.status=?", status).
		Select("orders.*, addresses.province as sender_province, addresses.city as sender_city, addresses.district as sender_district, addresses.street as sender_street, addresses.detail as sender_detail, addresses.receiver as sender_receiver, addresses.phone as sender_phone, addresses.latitude as sender_latitude, addresses.longitude as sender_longitude")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	currentPage := pageInfo.CurrentPage
	if currentPage < 1 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageInfo.PerPage
	err := query.Offset(offset).Limit(pageInfo.PerPage).Find(&orderList).Error
	if err != nil {
		return nil, total, err
	}
	return orderList, total, nil
}

// 骑手接单
func (r *ordersRepository) AcceptOrder(ctx context.Context, orderNo string, deliveryUserID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新订单状态为已接单待配送
		var order model.Order
		if err := tx.Where("order_no=?", orderNo).First(&order).Error; err != nil {
			return err
		}
		order.Status = consts.OrderBeforePickUp
		order.DeliveryUserID = &deliveryUserID
		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		//插入配送关系
		if err := tx.Create(&model.DeliveryAssign{
			OrderID:        order.ID,
			DeliveryUserID: deliveryUserID,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}
func (r *ordersRepository) PickupOrder(ctx context.Context, orderNo string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态为已取件待配送
		return tx.Model(&model.Order{}).Where("order_no=?", orderNo).Update("status", consts.OrderPickup).Error
	})
}

// 开始配送
func (r *ordersRepository) StartDelivery(ctx context.Context, orderNo string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态为配送中
		return tx.Model(&model.Order{}).Where("order_no=?", orderNo).Update("status", consts.OrderDelivering).Error
	})
}

// 上传位置
func (r *ordersRepository) UploadLocation(ctx context.Context, orderID uint, deliveryUserID uint, longitude, latitude float64) error {
	// 创建位置轨迹记录
	locationTrace := model.LocationTrace{
		OrderID:        orderID,
		DeliveryUserID: deliveryUserID,
		Longitude:      longitude,
		Latitude:       latitude,
		UploadTime:     time.Now(),
	}
	return r.db.Create(&locationTrace).Error
}

// 完成订单
func (r *ordersRepository) CompleteOrder(ctx context.Context, orderNo string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新订单状态为已完成
		var order model.Order
		if err := tx.Where("order_no=?", orderNo).First(&order).Error; err != nil {
			return err
		}
		order.Status = consts.OrderFinished
		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		// 2. 删除配送关系
		if err := tx.Where("order_id=?", order.ID).Delete(&model.DeliveryAssign{}).Error; err != nil {
			return err
		}
		return nil
	})
}
