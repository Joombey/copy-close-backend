package di

import (
	"dev.farukh/copy-close/config"
	"dev.farukh/copy-close/repos"
)

type Component struct {
	UserRepo repos.UserRepo
}

var component *Component

func init() {
	repo, err := repos.New(config.GetDSN())
	if err != nil {
		panic(err)
	}

	component = &Component{
		UserRepo: repo,
	}
}

func GetComponent() Component {
	return *component
}