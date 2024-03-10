package db_models

type Role struct {
	ID      uint  `gorm:"primarykey;autoIncrement"`
	CanSell *bool `gorm:"default:false"`
	CanBan  *bool `gorm:"default:false"`
	CanBuy  *bool `gorm:"default:true"`
	// for utility puprpose 
	Users []User
}
