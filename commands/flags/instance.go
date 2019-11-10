package flags

import "github.com/urfave/cli"

var InstanceServiceIdFlags = append(DurationFlags,
	cli.StringFlag{
		Required: true,
		Name:     "service",
		Usage:    "query service `ID`",
	})
