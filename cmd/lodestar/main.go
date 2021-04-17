package main

import (
	"github.com/lodestar-cli/lodestar/internal/cli/app"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var yamlKeys string
	var username string
	var token string
	var srcEnv string
	var name string
	var appConfigPath string
	var destEnv string
	var environment string
	var outputState bool

	a := &cli.App{
		Name: "lodestar",
		Version: "0.1.1",
		Usage: "Help guide your applications through their environments",
		Commands: []*cli.Command{
			{
				Name:        "app",
				Usage:       "Manage application images",
				Subcommands: []*cli.Command{
					{
						Name:  "push",
						Usage: "Push yaml key values to an environment",
						UsageText: "In order to push a tag to an environment, either a name for an App configured in ~/.lodestar\n\t"+
							       " needs to be provided with --name, or a path to an App needs to be provided with --config-path.\n\t"+
							       " Lodestar will then be able to find the App and pass the tag to the correct environment.",
						Flags: []cli.Flag {
							&cli.StringFlag{
								Name: "username",
								Hidden: true,
								Aliases: []string{"u"},
								Usage: "`username` for the version control account that can access the repository",
								Required: true,
								Destination: &username,
								EnvVars: []string{"GIT_USER"},
							},
							&cli.StringFlag{
								Name: "token",
								Hidden: true,
								Aliases: []string{"t"},
								Usage: "`token` for the version control account that can access the repository",
								Required: true,
								Destination: &token,
								EnvVars: []string{"GIT_TOKEN"},
							},
							&cli.StringFlag{
								Name: "name",
								Usage: "the `name` of an app",
								Destination: &name,
							},
							&cli.StringFlag{
								Name: "config-path",
								Usage: "the `path` to the app configuration file",
								Destination: &appConfigPath,
							},
							&cli.StringFlag{
								Name: "environment",
								Aliases: []string{"env"},
								Usage: "the `environment` the new yaml keys will be pushed to",
								Required: true,
								Destination: &environment,
							},
							&cli.StringFlag{
								Name: "yaml-keys",
								Usage: "a  comma separated `\"key=value\"` string of yaml keys to update",
								Destination: &yamlKeys,
								EnvVars: []string{"YAML_KEYS"},
							},
							&cli.BoolFlag{
								Name: "output-state",
								Usage: "will create a local yaml file of the updated app state when set",
								Destination: &outputState,
							},
						},
						Action: func(c *cli.Context) error {
							p, err := app.NewPush(username,token,name, appConfigPath, environment, yamlKeys)
							if err != nil {
								return err
							}
							err = p.Execute()
							if err != nil {
								return err
							}

							err = p.Output(outputState)
							if err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "promote",
						Usage: "Promote an image tag to the next environment",
						UsageText: "In order to promote an environment's tag, either a name for an App configured in ~/.lodestar\n\t"+
							" needs to be provided with --name, or a path to an a needs to be provided with --config-path.\n\t"+
							" Lodestar will then be able to find the App and pass the tag to the correct environment.",
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
								Name: "name",
								Usage: "the `name` of an a",
								Destination: &name,
							},
							&cli.StringFlag{
								Name: "config-path",
								Usage: "the `path` to the a configuration file",
								Destination: &appConfigPath,
							},
							&cli.StringFlag{
								Name: "src-env",
								Usage: "the `name` of the source environment",
								Required: true,
								Destination: &srcEnv,
							},
							&cli.StringFlag{
								Name: "dest-env",
								Usage: "the `name` of the destination",
								Required: true,
								Destination: &destEnv,
							},
							&cli.BoolFlag{
								Name: "output-state",
								Usage: "will create a local yaml file of the updated a state when set",
								Destination: &outputState,
							},
						},
						Action: func(c *cli.Context) error {
							p, err := app.NewPromote(username, token, name, appConfigPath, srcEnv, destEnv)
							if err != nil {
								return err
							}
							err = p.Execute()
							if err != nil {
								return err
							}

							err = p.Output(outputState)
							if err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "list",
						Usage: "List current context Apps",
						UsageText: "Will provide all the Apps within the current context as well as a description of the a.\n\t"+
							" App names and descriptions come directly from the appInfo block in their respective App configuration file.",
						Action: func(c *cli.Context) error {
							l, err := app.NewList()
							if err != nil{
								return err
							}
							l.Execute()
							return err
						},
					},
					{
						Name:  "show",
						Usage: "Prints the configuration file for the specified App",
						Flags: []cli.Flag{
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
								Name:        "name",
								Usage:       "the `name` of the a",
								Destination: &name,
							},
							&cli.StringFlag{
								Name: "config-path",
								Usage: "the `path` to the a configuration file",
								Destination: &appConfigPath,
							},
						},
						Action: func(c *cli.Context) error {
							s,err := app.NewShow(username,token,name,appConfigPath)
							if err != nil{
								return err
							}
							s.Execute()
							return nil
						},
					},
				},
			},
		},
	}

	err := a.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}