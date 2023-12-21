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
					index, e := strconv.Atoi(val)

					if e != nil {
						pterm.Error.Printf("Invalid index!\n")
						os.Exit(-1)
					}

					manga_id := mangas[index-1].Attrs()["href"]

					// spinner, _ = pterm.DefaultSpinner.Start()
					url := fmt.Sprintf("https://mangapill.com%s", manga_id)

					chapters := scrapeChapters(url)
					// spinner.Stop()
					pterm.Info.Printfln("Found %d chapters!", len(chapters))

					val, _ = pterm.DefaultInteractiveTextInput.Show("Select chapter: ")
					// Chapter order is reversed.
					index, e = strconv.Atoi(val)

					if e != nil {
						pterm.Error.Printf("Invalid index!\n")
						os.Exit(-1)
					}

					// Chapters' orders are reversed.
					pos := len(chapters) - index
					url = fmt.Sprintf("https://mangapill.com%s", chapters[pos])
					images := scrapeImages(url)

					downloadImages(images)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
