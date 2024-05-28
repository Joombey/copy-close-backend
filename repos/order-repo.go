package repos

import (
	apimodels "dev.farukh/copy-close/models/api_models"
	"dev.farukh/copy-close/models/db_models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func NewOrderRepo(dsn string) OrderRepo {
	db, _ := openConnection(dsn)

	return &OrderRepoImpl{
		db: db.Debug(),
	}
}

type OrderRepo interface {
	CreateOrder(request apimodels.OrderCreationRequests) error
	GetOrderInfo(userID uuid.UUID) OrderList
	UpdateOrderState(id uuid.UUID, newState int)
}

type OrderRepoImpl struct {
	db *gorm.DB
}

func (r *OrderRepoImpl) CreateOrder(request apimodels.OrderCreationRequests) error {
	order := db_models.Order{
		UserID:  uuid.FromStringOrNil(request.UserID),
		Comment: request.Comment,
	}
	r.db.Create(&order)

	go func() {
		for _, id := range request.Attachments {
			r.db.
				Model(&db_models.Document{}).
				Where("id = ?", uuid.FromStringOrNil(id)).
				Update("order_id", order.ID)
		}
	}()

	for _, servicePair := range request.Services {
		serviceID := uuid.FromStringOrNil(servicePair.First)
		amount := servicePair.Second

		var service db_models.Service
		r.db.Model(&db_models.Service{}).Where("id = ?", serviceID).Find(&service)

		orderService := db_models.OrderService{
			OrderID:   order.ID,
			ServiceID: service.ID,
			Price:     service.Price,
			Amount:    uint(amount),
			Title:     service.Title,
		}
		r.db.Model(&db_models.OrderService{}).Create(&orderService)
	}

	return nil
}

func (r *OrderRepoImpl) GetOrderInfo(userID uuid.UUID) OrderList {
	return OrderList{
		MyOrders: r.getOrdersForUser(userID),
		ToMe:     r.getOrdersForSeller(userID),
	}
}

func (r *OrderRepoImpl) UpdateOrderState(id uuid.UUID, newState int) {
	r.db.Model(&db_models.Order{}).Where("id = ?", id).UpdateColumn("state", newState)
}

func (r *OrderRepoImpl) getOrdersForSeller(userID uuid.UUID) []Orders {
	toMe := make([]Orders, 0)
	var user db_models.User
	r.db.Model(&db_models.User{}).Where("id = ?", userID).Find(&user)
	if user.RoleID == sellerRole.ID {
		var userServices []db_models.Service
		r.db.Model(&db_models.Service{}).Where("user_id = ?", userID).Find(&userServices)

		var userServicesIds []uuid.UUID
		for _, service := range userServices {
			userServicesIds = append(userServicesIds, service.ID)
		}

		var ordersWithServices []db_models.OrderService
		r.db.
			Model(&db_models.OrderService{}).
			Where("service_id in (?)", userServicesIds).
			Find(&ordersWithServices)

		var orderIds []uuid.UUID
		for _, orderService := range ordersWithServices {
			orderIds = append(orderIds, orderService.OrderID)
		}

		var orders_ []db_models.Order
		r.db.
			Model(&db_models.Order{}).
			Where("id in (?)", orderIds).
			Find(&orders_)

		orderUserMap := make(map[uuid.UUID]bool)
		for _, v := range orders_ {
			orderUserMap[v.UserID] = false
		}
		for userID := range orderUserMap {
			toMe = append(toMe, r.getOrdersForUser(userID)...)
		}
	}
	return toMe
}

func (r *OrderRepoImpl) getOrdersForUser(userID uuid.UUID) []Orders {
	var orders []db_models.Order
	r.db.Where("user_id = ?", userID).Find(&orders)

	result := make([]Orders, 0)
	for _, order := range orders {
		var orderService []db_models.OrderService
		r.db.
			Model(&db_models.OrderService{}).
			Where("order_id", order.ID).
			Find(&orderService)

		orderServiceSlice := make([]OrderService, 0)
		for _, v := range orderService {
			orderServiceSlice = append(orderServiceSlice, OrderService{
				ID:     v.ServiceID,
				Price:  int(v.Price),
				Amount: int(v.Amount),
				Title:  v.Title,
			})
		}

		var docs []db_models.Document
		r.db.Where("order_id", order.ID).Find(&docs)
		docsIds := make([]Attachment, 0)
		for _, v := range docs {
			docsIds = append(docsIds, Attachment{ID: v.ID, Name: v.Name})
		}

		var seller db_models.Service
		r.db.
			Model(&db_models.Service{}).
			Where("id = ?", orderServiceSlice[0].ID).
			Find(&seller)

		result = append(
			result,
			Orders{
				OrderID:     order.ID,
				Comment:     *order.Comment,
				SellerID:    *seller.UserID,
				State:       order.State,
				UserID:      userID,
				Services:    orderServiceSlice,
				Attachments: docsIds,
			},
		)
	}
	return result
}

type OrderService struct {
	ID     uuid.UUID `json:"id"`
	Price  int       `json:"price"`
	Amount int       `json:"amount"`
	Title  string    `json:"title"`
}

type Orders struct {
	OrderID     uuid.UUID      `json:"order_id"`
	Comment     string         `json:"comment"`
	SellerID    uuid.UUID      `json:"seller_id"`
	UserID      uuid.UUID      `json:"user_id"`
	State       int            `json:"state"`
	Services    []OrderService `json:"services"`
	Attachments []Attachment   `json:"attachments"`
}

type OrderList struct {
	MyOrders []Orders `json:"my_orders"`
	ToMe     []Orders `json:"to_me"`
}

type Attachment struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
