package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/shanna/git-vendor/vendor"
	"log"
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
	vc, err := vendor.Open(".")
	if err != nil {
		log.Fatalf(err.Error())
	}

	if len(vc.Repositories) == 0 {
		log.Fatalf("no " + vendor.Filename + ", nothing to vendor")
	}
	for _, repository := range vc.Repositories {
		log.Printf("vendor " + repository.Name)
		if err := repository.Vendor(); err != nil {
			log.Fatalf(fmt.Sprintf("failed to vendor %s: %s", repository.Name, err.Error()))
		}
	}
}
