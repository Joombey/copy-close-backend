package repos

import (
	api "dev.farukh/copy-close/models/api_models"
	dbModels "dev.farukh/copy-close/models/db_models"
	errs "dev.farukh/copy-close/models/errs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var userRole = dbModels.Role{
	CanSell: boolPtr(false),
	CanBan:  boolPtr(false),
	CanBuy:  boolPtr(true),
}

var sellerRole = dbModels.Role{
	CanSell: boolPtr(true),
	CanBan:  boolPtr(false),
	CanBuy:  boolPtr(true),
}

var adminRole = dbModels.Role{
	CanSell: boolPtr(false),
	CanBan:  boolPtr(true),
	CanBuy:  boolPtr(false),
}

type UserRepo interface {
	CreateUser(api.SignUpRequest) error
}

type UserRepoImpl struct {
	db *gorm.DB
}

func (repo UserRepoImpl) CreateUser(signUp api.SignUpRequest) error {
	user := &dbModels.User{
		Login:      signUp.Login,
		FirstName:  signUp.Name,
		SecondName: signUp.SecondName,
		Password:   signUp.Password,
		Address:    signUp.Address,
		RoleID: func() uint {
			if signUp.IsSeller == boolPtr(true) {
				return sellerRole.ID
			} else {
				return userRole.ID
			}
		}(),
	}
	var exists bool
	repo.db.Raw("SELECT EXISTS(SELECT * FROM users WHERE login = ? and password = ?)", signUp.Login, signUp.Password).Scan(&exists)
	if !exists {
		err := repo.db.Create(user)
		if err.Error != nil {
			return err.Error
		}
		return nil
	} else {
		return errs.ErrUserExists
	}
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

	createIfNotExists(db, &userRole)
	createIfNotExists(db, &sellerRole)
	createIfNotExists(db, &adminRole)

	repo := &UserRepoImpl{
		db: db.Debug(),
	}
	return repo, nil
}

func openConnection(dsn string) (*gorm.DB, error) {
	mysqlConfig := mysql.Open(dsn)
	return gorm.Open(mysqlConfig, &gorm.Config{})
}

func setupDb(db *gorm.DB) error {
	return db.AutoMigrate(
		&dbModels.Role{},
		&dbModels.User{},
	)
}

func createIfNotExists(db *gorm.DB, role *dbModels.Role) {
	var tmpRole dbModels.Role
	db.Debug().Where(role).Find(&tmpRole)
	if tmpRole.CanBan == nil {
		db.Create(role)
	} else {
		role.ID = tmpRole.ID
	}
}

func boolPtr(b bool) *bool {
	return &b
}
