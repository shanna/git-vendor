package vendor

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"strings"
)

type Vendor struct {
	Config       string
	Repositories []*Repository
}

type Repository struct {
	Name   string
	Url    string
	Path   string
	Branch string
}

func Load(config string) (*Vendor, error) {
	v := &Vendor{
		Config:       config,
		Repositories: make([]*Repository, 0),
	}

	parent, err := git.NewConfig()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating git config: %s", err.Error()))
	}

	vendor, err := git.OpenOndisk(parent, config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading config: %s", err.Error()))
	}

	it, _ := vendor.NewIteratorGlob(`vendor\..*?\.url`)
	for {
		entry, _ := it.Next()
		if entry == nil {
			break
		}

		parts := strings.Split(entry.Name, ".")
		repository := &Repository{
			Name: parts[1],
			Url:  entry.Value,
		}

		if path, err := vendor.LookupString(fmt.Sprintf("vendor.%s.path", repository.Name)); err == nil {
			repository.Path = path
		}
		if branch, err := vendor.LookupString(fmt.Sprintf("vendor.%s.branch", repository.Name)); err == nil {
			repository.Branch = branch
		}

		v.Repositories = append(v.Repositories, repository)
	}

	return v, nil
}
