package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "caldera-admin-launcher",
		Action: walkGitRepositoryCommits,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     ArgRepositoryLocation,
				Aliases:  []string{"r"},
				Usage:    "Location of the git repository to scan",
				Required: true,
			},
			&cli.StringFlag{
				Name:     ArgFromDate,
				Aliases:  []string{"d"},
				Usage:    "If set, only show commits from this date forward",
				Required: false,
			},
			&cli.StringFlag{
				Name:     ArgOutputFile,
				Aliases:  []string{"o"},
				Usage:    "Output commit history to this file",
				Required: false,
				Value:    "output.json",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
