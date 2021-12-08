package main

import (
	"log"
	"os"
	"sort"

	"github.com/amorist/assimp-go/assimp"
	"github.com/urfave/cli/v2"
)

func main() {
	var input, output, format string
	app := &cli.App{
		Name:        "assimp-go",
		Description: "A simple command line tool to convert models to glTF 2.0",
		Version:     "0.0.1",
		Commands: []*cli.Command{
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "export",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input",
						Aliases:     []string{"i"},
						Value:       "",
						Usage:       "input file",
						Destination: &input,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Value:       "",
						Usage:       "output file",
						Destination: &output,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "format",
						Aliases:     []string{"f"},
						Value:       "",
						Usage:       "output format",
						Destination: &format,
					},
				},
				Action: func(c *cli.Context) error {
					if len(format) == 0 {
						format = "obj"
					}
					err := assimp.Export(input, output, format)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "import",
				Aliases: []string{"im"},
				Usage:   "import",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input",
						Aliases:     []string{"i"},
						Value:       "",
						Usage:       "input file",
						Destination: &input,
						Required:    true,
					},
				},
				Action: func(c *cli.Context) error {
					_, err := assimp.Import(input)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
