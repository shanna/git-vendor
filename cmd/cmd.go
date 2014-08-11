package cmd

import (
	"github.com/codegangsta/cli"
	git "github.com/libgit2/git2go"
	"log"
	"os"
)

const vendorFilename string = ".gitvendor"

var repoDirectory string

func setup(ctx *cli.Context) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %s", err.Error())
	}

	// TODO: Search upwards in case of nested projects.
	repoDirectory, err = git.Discover(wd, true, []string{"/"})
	if err != nil {
		log.Fatal("failed to get git working directory: %s", err.Error())
	}
}
