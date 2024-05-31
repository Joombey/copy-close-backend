package repos

import (
	"dev.farukh/copy-close/models/db_models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type AdminRepo interface {
	GetBlocklist() []Report
	SetSolution(reportID uuid.UUID, state int)
}

type AdminRepoImpl struct {
	db *gorm.DB
}

func NewAdminRepo(dsn string) AdminRepo {
	db, _ := openConnection(dsn)
	return &AdminRepoImpl{
		db: db.Debug(),
	}
}

func (r *AdminRepoImpl) GetBlocklist() []Report {
	var reports []db_models.Report
	r.db.Where("solution = ?", 0).Order("created_at asc").Find(&reports)

	result := make([]Report, 0, len(reports))
	for _, report := range reports {
		var order *db_models.Order
		r.db.Where("id = ?", report.OrderID).Find(&order)

		var user *db_models.User
		r.db.Where("id = ?", order.UserID).Find(&user)

		var orderService *db_models.OrderService
		r.db.Where("order_id = ?", order.ID).Find(&orderService)

		var service *db_models.Service
		r.db.Where("id = ?", orderService.ServiceID).Find(&service)

		var seller *db_models.User
		r.db.Where("id = ?", service.UserID).Find(&seller)

		result = append(result, Report{
			ReportID:      report.ID,
			OrderID:       order.ID,
			UserID:        user.ID,
			SellerID:      seller.ID,
			ReportMessage: report.Message,
			OrderMessage:  order.Comment,
			ReportDate:    report.CreatedAt.String(),
		})
	}
	return result
}

func (r *AdminRepoImpl) SetSolution(reportID uuid.UUID, state int) {
	r.db.
		Model(&db_models.Report{}).
		Where("id = ?", reportID).
		Update("solution", state)
}

type Report struct {
	ReportID      uuid.UUID `json:"report_id"`
	OrderID       uuid.UUID `json:"order_id"`
	UserID        uuid.UUID `json:"user_id"`
	SellerID      uuid.UUID `json:"seller_id"`
	ReportMessage string    `json:"report_message"`
	OrderMessage  *string   `json:"order_message"`
	ReportDate    string    `json:"report_date"`
}
