package flags

import "github.com/urfave/cli"

var SearchRegexFlags = []cli.Flag{
	cli.StringFlag{
		Name:     "regex",
		Required: true,
		Usage:    "search `Regex`",
	},
}
