package repos

import (
	"dev.farukh/copy-close/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(fName, sName, passwordHash string)
}

type UserRepoImpl struct {
	users *gorm.DB
	db    *gorm.DB
}

func (repo UserRepoImpl) CreateUser(fName, sName, passwordHash string) {
	user := &models.User{
		FirstName: fName,
		SecondName: sName,
		Password: passwordHash,
		Role: 
	}
	repo.users.Create()
}

func New(dsn string) (*UserRepoImpl, error) {
	db, err := openConnection(dsn)
	if err != nil {
		return nil, err
	}

	err = setupDb(db)
	if err != nil {
		return nil, err
	}

	repo := &UserRepoImpl{
		db:    db,
		users: db.Table("users"),
	}
	return repo, nil
}

func openConnection(dsn string) (*gorm.DB, error) {
	mysqlConfig := mysql.Open(dsn)
	return gorm.Open(mysqlConfig, &gorm.Config{})
}

func setupDb(db *gorm.DB) error {
	return db.AutoMigrate(
		models.User{},
	)
}
