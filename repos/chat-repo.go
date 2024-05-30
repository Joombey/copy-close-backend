package repos

import (
	"dev.farukh/copy-close/models/db_models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ChatMessage struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`

	MessageID uuid.UUID `json:"message_id"`
	Text      string    `json:"text"`
}

func NewChatRepo(dsn string) ChatRepo {
	db, _ := openConnection(dsn)
	return &ChatRepoImpl{
		db: db.Debug(),
	}
}

type ChatRepo interface {
	GetMessagesForOrder(orderID uuid.UUID) []ChatMessage
	CreateMessage(orderID, userID uuid.UUID, text string) error
}

type ChatRepoImpl struct {
	db *gorm.DB
}

func (repo *ChatRepoImpl) GetMessagesForOrder(orderID uuid.UUID) []ChatMessage {
	var messages []db_models.Message
	rows := repo.db.
		Model(&db_models.Message{}).
		Where("order_id = ?", orderID).
		Order("created_at asc").
		Find(&messages).RowsAffected

	if rows == 0 {
		return []ChatMessage{}
	}

	var result []ChatMessage
	userMap := make(map[uuid.UUID]*db_models.User)
	for _, message := range messages {
		user, ok := userMap[message.UserID]
		if !ok {
			repo.db.Where("id = ?", message.UserID).Find(&user)
			userMap[message.UserID] = user
		}

		result = append(result, ChatMessage{
			UserID:   message.UserID,
			UserName: userMap[message.UserID].FirstName,

			MessageID: message.ID,
			Text:      message.Text,
		})
	}

	return result
}

func (repo *ChatRepoImpl) CreateMessage(orderID, userID uuid.UUID, text string) error {
	return repo.db.Create(&db_models.Message{
		OrderID: orderID,
		UserID:  userID,
		Text:    text,
	}).Error
}
