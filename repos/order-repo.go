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
