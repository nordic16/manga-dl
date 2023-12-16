package main

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/pterm/pterm"
	cli "github.com/urfave/cli/v2"
)

/* Contains a bunch of useful functions to handle cli commands. */

// TODO: split this code into (maybe) a few functions to make it better.
func start() {
	app := &cli.App{
		Name:  "manga-dl",
		Usage: "Epic terminal manga reader lmfao",
		Commands: []*cli.Command{
			{
				Name:  "search",
				Usage: "searches for a given manga",
				Action: func(cCtx *cli.Context) error {
					var manga string = cCtx.Args().First()
					spinner, _ := pterm.DefaultSpinner.Start()

					mangas := scrapeMangas(manga)
					listMangas(mangas)

					spinner.Stop()

					val, _ := pterm.DefaultInteractiveTextInput.Show("Select index: ")
					index, _ := strconv.Atoi(val)

					url := fmt.Sprintf("https://mangareader.to%s", getUrl(mangas[index-1]))

					fmt.Println(url)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
