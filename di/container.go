package di

import (
	"fmt"

	"dev.farukh/copy-close/config"
	"dev.farukh/copy-close/repos"
)

type Component struct {
	UserRepo repos.UserRepo
}

func New() Component {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	
	repo, err := repos.New(cfg.DSN())
	if err != nil {
		panic(err)
	}

	return Component{
		UserRepo: &repo,
	}
}