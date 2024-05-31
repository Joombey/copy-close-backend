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

type UserRepo interface {
	RegisterUser(api.RegisterRequest, bool) (repoModels.RegisterResult, error)
	LogInUser(api.LogInRequest) (uuid.UUID, error)
	GetUser(string, string, string) (repoModels.UserInfoResult, error)
	GetUserInternal(string) (repoModels.UserInfoResult, error)
	GetSellers() []repoModels.UserInfoResult
	EditProfile(api.EditProfileRequest, string) error
	CheckTokenValid(string, string) bool
	DeleteUser(string)
}

type UserRepoImpl struct {
	db *gorm.DB
}

func (repo UserRepoImpl) RegisterUser(registerData api.RegisterRequest, admin bool) (repoModels.RegisterResult, error) {
	if !repo.userExists(registerData.Login, registerData.Password) {
		addressID := repo.createAddress(registerData.Address)

		var role dbModels.Role
		if admin {
			role = adminRole
		} else {
			role = repo.getRole(registerData.IsSeller)
		}

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

func (repo UserRepoImpl) LogInUser(signIn api.LogInRequest) (authToken uuid.UUID, err error) {
	user := &dbModels.User{}
	err = repo.db.Where("login = ? AND password = ?", signIn.Login, signIn.Password).First(user).Error
	if err != nil {
		return uuid.UUID{}, errs.ErrInvalidLoginOrPassword
	}
	return user.AuthToken, nil
}

func (repo UserRepoImpl) GetUser(login, authToken, id string) (repoModels.UserInfoResult, error) {
	if tokenValid := repo.CheckTokenValid(id, authToken); id != "" && !tokenValid {
		return repoModels.UserInfoResult{}, errs.ErrInvalidLoginOrAuthToken
	}

	isUUID := uuid.FromStringOrNil(login) == uuid.Nil

	var user dbModels.User
	var err error
	if !isUUID {
		if id == "" {
			err = repo.db.Where("id = ? AND auth_token = ?", login, authToken).First(&user).Error
		} else {
			err = repo.db.Where("id = ?", login).First(&user).Error
		}

	} else {
		if id == "" {
			err = repo.db.Where("login = ? AND auth_token = ?", login, authToken).First(&user).Error
		} else {
			err = repo.db.Where("login = ?", login).First(&user).Error
		}
	}

	if err != nil {
		return repoModels.UserInfoResult{}, err
	}

	var services []dbModels.Service
	repo.db.Where("user_id = ? AND deleted = 0", user.ID).Find(&services)

	var role dbModels.Role
	err = repo.db.Where("id = ?", user.RoleID).First(&role).Error
	if err != nil {
		return repoModels.UserInfoResult{}, err
	}

	var address dbModels.Address
	err = repo.db.Where("id = ?", user.AddressID).First(&address).Error
	if err != nil {
		return repoModels.UserInfoResult{}, err
	}
	return repoModels.UserInfoResult{
		User:     user,
		Role:     role,
		Address:  address,
		Services: services,
	}, nil
}

func (repo UserRepoImpl) GetUserInternal(userID string) (repoModels.UserInfoResult, error) {
	var user dbModels.User
	err := repo.db.Where("id = ?", uuid.FromStringOrNil(userID)).First(&user).Error

	if err != nil {
		return repoModels.UserInfoResult{}, err
	}

	var services []dbModels.Service
	repo.db.Where("user_id = ? AND deleted = 0", user.ID).Find(&services)

	var role dbModels.Role
	err = repo.db.Where("id = ?", user.RoleID).First(&role).Error
	if err != nil {
		return repoModels.UserInfoResult{}, err
	}

	var address dbModels.Address
	err = repo.db.Where("id = ?", user.AddressID).First(&address).Error
	if err != nil {
		return repoModels.UserInfoResult{}, err
	}
	return repoModels.UserInfoResult{
		User:     user,
		Role:     role,
		Address:  address,
		Services: services,
	}, nil
}

func (repo UserRepoImpl) GetSellers() []repoModels.UserInfoResult {
	var sellers []dbModels.User
	repo.db.Model(
		&dbModels.User{},
	).Preload(
		"Users",
	).Where(
		"role_id = ?", sellerRole.ID,
	).Find(&sellers)

	var infos []repoModels.UserInfoResult
	for _, seller := range sellers {
		repoModel, _ := repo.GetUser(seller.Login, seller.AuthToken.String(), seller.ID.String())
		infos = append(infos, repoModel)
	}
	return infos
}

func (repo UserRepoImpl) CheckTokenValid(userID, authToken string) bool {
	tokenUUID, err := uuid.FromString(authToken)
	if err != nil {
		return false
	}

	var exists bool
	repo.db.Raw(
		"SELECT EXISTS (SELECT * FROM users WHERE id = ? AND auth_token = ?)",
		userID,
		tokenUUID,
	).Scan(&exists)

	return exists
}

func (repo UserRepoImpl) EditProfile(editProfileRequest api.EditProfileRequest, imageID string) error {
	var user dbModels.User
	repo.db.Where("id = ?", editProfileRequest.UserID).First(&user)

	if user.AuthToken.String() != editProfileRequest.AuthToken {
		return errs.ErrInvalidLoginOrAuthToken
	}

	user.FirstName = editProfileRequest.Name
	if imageID != "" {
		user.UserImage = uuid.FromStringOrNil(imageID)
	}

	repo.db.UpdateColumns(user)

	if user.RoleID != sellerRole.ID {
		return nil
	}

	for _, service := range editProfileRequest.Services {
		if repo.serviceExists(service) {
			var actucalService dbModels.Service
			repo.db.Where("id = ?", service.ID).First(&actucalService)
			actucalService.Title = service.Title
			actucalService.Price = service.Price

			repo.db.UpdateColumns(actucalService)
		} else {
			newService := dbModels.Service{
				UserID: &user.ID,
				Title:  service.Title,
				Price:  service.Price,
			}
			repo.db.Create(&newService)
		}
	}
	for _, serviceID := range editProfileRequest.ServicesToDelete {
		repo.db.Raw(
			"UPDATE services SET deleted = 1 WHERE id = ?",
			uuid.FromStringOrNil(serviceID),
		).Scan(nil)
	}
	return nil
}

func (repo UserRepoImpl) serviceExists(service dbModels.Service) bool {
	if service.ID == uuid.Nil {
		return false
	}

	var exists bool
	repo.db.Raw(
		"SELECT EXISTS (SELECT * FROM services WHERE id = ?)",
		service.ID,
	).Scan(&exists)

	return exists
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
	if *isSeller {
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

func (repo *UserRepoImpl) DeleteUser(userId string) {
	repo.db.Exec(
		"DELETE FROM users WHERE user_id = ?",
		uuid.FromStringOrNil(userId),
	)
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
	err := db.Debug().SetupJoinTable(
		&dbModels.Order{},
		"Services",
		&dbModels.OrderService{},
	)
	if err != nil {
		return err
	}

	return db.Debug().AutoMigrate(
		&dbModels.Role{},
		&dbModels.Address{},
		&dbModels.User{},
		&dbModels.Service{},
		&dbModels.Document{},
		&dbModels.Order{},
		&dbModels.OrderService{},
		&dbModels.Message{},
		&dbModels.Report{},
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
