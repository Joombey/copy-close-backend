package db_models

type Role struct {
	ID      uint  `gorm:"primarykey;autoIncrement" json:"id"`
	CanSell *bool `gorm:"default:false" json:"can_sell"`
	CanBan  *bool `gorm:"default:false" json:"can_ban"`
	CanBuy  *bool `gorm:"default:true" json:"can_buy"`
	// for utility puprpose 
	Users []User `json:"-"`
}
