package main

import (
	lodeApp "github.com/lodestar-cli/lodestar/internal/app"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var url string
	var tag string
	var username string
	var token string
	var srcPath string
	var destPath string

	app := &cli.App{
		Name: "lodestar",
		Version: "0.0.1",
		Usage: "Help guide your applications through their environments",
		Commands: []*cli.Command{
			{
				Name:        "app",
				Usage:       "Manage application images",
				Subcommands: []*cli.Command{
					{
						Name:  "push",
						Usage: "push a new image tag to an environment",
						UsageText: "Locates a specified environment configuration manifest within a repository, and updates the application's tag for that environment",
						Flags: []cli.Flag {
							&cli.StringFlag{
								Name: "username",
								Hidden: true,
								Usage: "`username` for the version control account that can access the repository",
								Required: true,
								Destination: &username,
								EnvVars: []string{"GIT_USER"},
							},
							&cli.StringFlag{
								Name: "token",
								Hidden: true,
								Usage: "`token` for the version control account that can access the repository",
								Required: true,
								Destination: &token,
								EnvVars: []string{"GIT_TOKEN"},
							},
							&cli.StringFlag{
								Name: "repo",
								Usage: "`url` of the app's manifest repository",
								Required: true,
								Destination: &url,
							},
							&cli.StringFlag{
								Name: "src-path",
								Usage: "the `path` to the environment's configuration file",
								Required: true,
								Destination: &srcPath,
							},
							&cli.StringFlag{
								Name: "tag",
								Usage: "the `tag` for the new image",
								Required: true,
								Destination: &tag,
							},
						},
						Action: func(c *cli.Context) error {
							err :=lodeApp.Push(username,token,url,srcPath,tag)

							return err
						},
					},
					{
						Name:  "promote",
						Usage: "promote an image tag to the next environment",
						UsageText: "Retrieves an application tag specified in a source environment configuration file, and promotes it to a destination configuration file",
						Flags: []cli.Flag {
							&cli.StringFlag{
								Name: "username",
								Usage: "`username` for the version control account that can access the repository",
								Hidden: true,
								Required: true,
								Destination: &username,
								EnvVars: []string{"GIT_USER"},
							},
							&cli.StringFlag{
								Name: "token",
								Usage: "`token` for the version control account that can access the repository",
								Hidden: true,
								Required: true,
								Destination: &token,
								EnvVars: []string{"GIT_TOKEN"},
							},
							&cli.StringFlag{
								Name: "repo",
								Usage: "`url` of the app's manifest repository",
								Required: true,
								Destination: &url,
							},
							&cli.StringFlag{
								Name: "src-path",
								Usage: "the `path` to the source environment's configuration file",
								Required: true,
								Destination: &srcPath,
							},
							&cli.StringFlag{
								Name: "dest-path",
								Usage: "the `path` to the destination environment's configuration file",
								Required: true,
								Destination: &destPath,
							},
						},
						Action: func(c *cli.Context) error {
							err := lodeApp.Promote(username,token,url,srcPath,destPath)
							return err
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}