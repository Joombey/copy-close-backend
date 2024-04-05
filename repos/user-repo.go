package repos

import (
	api "dev.farukh/copy-close/models/api_models"
	core "dev.farukh/copy-close/models/core_models"
	dbModels "dev.farukh/copy-close/models/db_models"
	errs "dev.farukh/copy-close/models/errs"
	repoModels "dev.farukh/copy-close/models/repo_models"
	uuid "github.com/satori/go.uuid"
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
	RegisterUser(api.RegisterRequest) (repoModels.RegisterResult, error)
	LogInUser(api.LogInRequest) error
}

type UserRepoImpl struct {
	db *gorm.DB
}

func (repo UserRepoImpl) RegisterUser(registerData api.RegisterRequest) (repoModels.RegisterResult, error) {
	if !repo.userExists(registerData.Login, registerData.Password) {
		addressID := repo.createAddress(registerData.Address)
		role := repo.getRole(registerData.IsSeller)

		id, authToken, imageID, err := repo.createUser(registerData, role.ID, addressID)
		if err != nil {
			return repoModels.RegisterResult{}, err
		}

		result := repoModels.RegisterResult{
			UserID:    id.String(),
			AddressID: addressID.String(),
			AuthToken: authToken.String(),
			UserImage: imageID.String(),
			Role:      role,
		}
		return result, nil
	} else {
		return repoModels.RegisterResult{}, errs.ErrUserExists
	}
}

func (repo UserRepoImpl) LogInUser(signIn api.LogInRequest) (err error) {
	var exists bool
	err = repo.db.Raw(
		"SELECT EXISTS(SELECT * FROM users WHERE login = ? AND password = ?)",
		signIn.Login,
		signIn.Password,
	).Scan(&exists).Error

	if !exists {
		return errs.ErrInvalidLoginOrPassword
	}
	return
}

func (repo UserRepoImpl) userExists(login, password string) bool {
	var exists bool
	repo.db.Raw(
		"SELECT EXISTS(SELECT * FROM users WHERE login = ? and password = ?)",
		login,
		password,
	).Scan(&exists)
	return exists
}

func (repo UserRepoImpl) getRole(isSeller *bool) dbModels.Role {
	if isSeller == boolPtr(true) {
		return sellerRole
	} else {
		return userRole
	}
}

func (repo UserRepoImpl) createUser(
	registerData api.RegisterRequest,
	roleID uint,
	addressID uuid.UUID,
) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	imageID := uuid.NewV4()
	authToken := uuid.NewV4()
	user := &dbModels.User{
		Login:     registerData.Login,
		FirstName: registerData.Name,
		Password:  registerData.Password,
		AddressID: addressID,
		RoleID:    roleID,
		UserImage: imageID,
		AuthToken: authToken,
	}

	err := repo.db.Create(user)
	return user.ID, authToken, imageID, err.Error
}

func (repo UserRepoImpl) createAddress(addressCore core.Address) uuid.UUID {
	address := &dbModels.Address{
		AddressName: addressCore.Address,
		Lat:         addressCore.Lat,
		Lon:         addressCore.Lon,
	}
	repo.db.Create(&address)
	return address.ID
}

func (repo *UserRepoImpl) ClearAll() {
	repo.db.Exec("DELETE FROM users")
	repo.db.Exec("DELETE FROM addresses")
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
		&dbModels.Address{},
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
