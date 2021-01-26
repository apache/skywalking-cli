package event

import (
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:  "event",
	Usage: "Event related sub-command",
	Subcommands: []cli.Command{
		reportCommand,
	},
}
