package di

import (
	"dev.farukh/copy-close/config"
	"dev.farukh/copy-close/repos"
)

type Component struct {
	UserRepo  repos.UserRepo
	FileRepo  repos.FileRepo
	OrderRepo repos.OrderRepo
	ChatRepo  repos.ChatRepo
}

var component *Component

func init() {
	dsn := config.GetDSN()
	repo, err := repos.New(dsn)
	if err != nil {
		panic(err)
	}

	component = &Component{
		UserRepo:  repo,
		FileRepo:  repos.NewFileRepo(dsn),
		OrderRepo: repos.NewOrderRepo(dsn),
		ChatRepo:  repos.NewChatRepo(dsn),
	}
}

func GetComponent() Component {
	return *component
}
