// git-vendor
package main

import (
	"github.com/codegangsta/cli"
	"github.com/shanna/git-vendor/cmd"
	"log"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	log.SetFlags(0)

	app := cli.NewApp()
	app.Name = "git-vendor"
	app.Usage = "Relaxed Git Submodules"
	app.Version = "0.1" // TODO:
	app.Commands = []cli.Command{
		// cmd.CmdAdd,
		// cmd.CmdRemove,
		cmd.CmdInstall,
		// cmd.CmdUpdate,
	}
	app.Run(os.Args)
}
