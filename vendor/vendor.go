package vendor

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"os"
	"os/user"
	"path"
	"strings"
)

var Filename string = ".gitvendor"

type Vendor struct {
	config       *git.Config
	Repositories []*Repository
}

type Repository struct {
	Name   string
	Url    string
	Path   string
	Branch string
}

func discover(path string) (*git.Config, error) {
	if path == "" || path == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, errors.New("failed to get current work directory: " + err.Error())
		}
		path = cwd
	}

	path, err := git.Discover(path, true, []string{"/"})
	if err != nil {
		return nil, errors.New("failed to get git work directory: " + err.Error())
	}

	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, errors.New("failed reading git repository: " + err.Error())
	}

	if err := os.Chdir(repo.Workdir()); err != nil {
		return nil, errors.New("failed to change work directory: " + err.Error())
	}

	parent, _ := git.NewConfig()
	config, err := git.OpenOndisk(parent, Filename)
	if err != nil {
		return nil, errors.New("failed to read config: " + err.Error())
	}
	return config, nil
}

func repositories(config *git.Config) ([]*Repository, error) {
	repos := make([]*Repository, 0)
	it, _ := config.NewIteratorGlob(`vendor\..*?\.url`)
	for {
		entry, _ := it.Next()
		if entry == nil {
			break
		}

		parts := strings.Split(entry.Name, ".")
		repo := &Repository{
			Name:   parts[1],
			Url:    entry.Value,
			Path:   parts[1],
			Branch: "master",
		}
		if repoPath, err := config.LookupString(fmt.Sprintf("vendor.%s.path", repo.Name)); err == nil {
			repo.Path = repoPath
		}
		if repoBranch, err := config.LookupString(fmt.Sprintf("vendor.%s.branch", repo.Name)); err == nil {
			repo.Branch = repoBranch
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func credentials(url string, username_from_url string, allowed_types git.CredType) (int, *git.Cred) {
	if allowed_types&git.CredTypeSshKey > 0 {
		currentUser, _ := user.Current()
		rc, cred := git.NewCredSshKey(
			username_from_url,
			path.Join(currentUser.HomeDir, ".ssh", "id_rsa.pub"),
			path.Join(currentUser.HomeDir, ".ssh", "id_rsa"),
			"",
		)
		return rc, &cred
	}
	rc, cred := git.NewCredDefault()
	return rc, &cred
}

func Open(path string) (*Vendor, error) {
	config, err := discover(path)
	if err != nil {
		return nil, err
	}

	repos, err := repositories(config)
	if err != nil {
		return nil, err
	}

	return &Vendor{
		config:       config,
		Repositories: repos,
	}, nil
}

func (r Repository) Vendor() error {
	if stat, err := os.Stat(r.Path); os.IsNotExist(err) {
		if err = os.MkdirAll(r.Path, 0776); err != nil {
			return err
		}

		options := &git.CloneOptions{
			RemoteCallbacks: &git.RemoteCallbacks{
				CredentialsCallback: credentials,
			},
		} // TODO: Progress

		if _, err := git.Clone(r.Url, r.Path, options); err != nil {
			return err
		}
	} else if !stat.IsDir() {
		return errors.New(fmt.Sprintf("%s exists, but is a file", r.Path))
	}

	repository, err := git.OpenRepository(r.Path)
	if err != nil {
		fmt.Printf("err:%+V\n", err)
		return err
	}

	remote, err := repository.LoadRemote("origin")
	if err != nil {
		return err
	}

	if err := remote.Fetch(nil, ""); err != nil {
		return err
	}

  reference, err := repository.DwimReference(r.Branch)
  if err != nil {
    return err
  }

  switch reference.Type() {
  case git.ReferenceSymbolic:
    return errors.New(fmt.Sprintf("symbolic references not supported yet: %s", r.Branch))
  case git.ReferenceOid:
    var signature git.Signature
    if err := repository.SetHead(reference.Name(), &signature, ""); err != nil {
      return err
    }

    options := &git.CheckoutOpts{
      Strategy: git.CheckoutSafe,
    } // CloneOptions but CheckoutOpts !?
    if err := repository.CheckoutHead(options); err != nil {
      return err
    }
  }
	return nil
}
