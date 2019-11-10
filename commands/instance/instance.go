package instance

import (
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:      "instance",
	ShortName: "i",
	Usage:     "Instance related sub-command",
	Subcommands: cli.Commands{
		ListCommand,
	},
}
