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

					val, _ := pterm.DefaultInteractiveTextInput.Show("Select manga: ")
					index, _ := strconv.Atoi(val)

					manga_id := mangas[index-1].Find("a").Attrs()["href"]

					url := fmt.Sprintf("https://mangareader.to%s", manga_id)
					chapters := totalChapters(url)

					pterm.Println("\nFound " + pterm.LightGreen(chapters) + " chapters!\n")
					val, _ = pterm.DefaultInteractiveTextInput.Show("Select chapter: ")
					index, _ = strconv.Atoi(val)

					// Final url
					pterm.Printfln("https://mangareader.to/read%s/en/chapter-%d", manga_id, index)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
