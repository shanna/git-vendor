package cmd

import (
	"github.com/codegangsta/cli"
	git "github.com/libgit2/git2go"
	"github.com/shanna/git-vendor/vendor"
	"log"
	"os"
)

var CmdInstall = cli.Command{
	Name:  "install",
	Usage: "install vendored git packages",
	Description: `

`,
	Action: runInstall,
	Flags:  []cli.Flag{}, // TODO: -r recursive. Read .gitvendor out of sub projects.
}

func runInstall(ctx *cli.Context) {
	setup(ctx)

	// TODO: The repo location stuff can probably be moved into vendor/* as well.

	repo, err := git.OpenRepository(repoDirectory)
	if err != nil {
		log.Fatalf("error reading git repository: %s", err.Error())
	}

	log.Printf("changing work directory to %s", repo.Workdir())
	if err := os.Chdir(repo.Workdir()); err != nil {
		log.Fatalf("failed to change work directory: ", err.Error())
	}

	vc, err := vendor.Load(vendorFilename)
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}

	for _, repository := range vc.Repositories {
		log.Printf("%+v", repository)
	}
}
