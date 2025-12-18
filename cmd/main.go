package main

import (
	"context"
	"github.com/alhamsya/bookcabin/cmd/public"
	"github.com/urfave/cli/v2"
	"os"
	"syscall"
)

func main() {

	cliApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run flight search service",
				Subcommands: []*cli.Command{
					{
						Name:  "public-api",
						Usage: "Run public REST API",
						Action: func(ctx *cli.Context) error {
							return public.RunRESTApp(ctx.Context, syscall.SIGINT, syscall.SIGTERM)
						},
					},
				},
			},
		},
		AllowExtFlags: true,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name: "cfg-credential",
			},
			&cli.IntFlag{
				Name: "cfg-static",
			},
		},
	}

	if err := cliApp.RunContext(context.Background(), os.Args); err != nil {
		panic(err.Error())
	}
}
