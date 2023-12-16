package main

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/anaskhan96/soup"
	"github.com/pterm/pterm"
	cli "github.com/urfave/cli/v2"
)

/* Contains a bunch of useful functions to handle cli commands. */

// TODO: split this code into (maybe) a few functions to make it better.
func start() {
	app := &cli.App{
		Name:  "manga-dl",
		Usage: "Epic Epic terminal manga reader lmfao",
		Commands: []*cli.Command{
			{
				Name:  "search",
				Usage: "searches for a given manga",
				Action: func(cCtx *cli.Context) error {
					var manga string = cCtx.Args().First()
					spinner, _ := pterm.DefaultSpinner.Start()

					resp, err := soup.Get("https://mangareader.to/search?keyword=" + manga)

					if err != nil {
						os.Exit(1)
					}

					// Actually handles the requests.
					doc := soup.HTMLParse(resp)
					mangas := doc.FindAll("h3", "class", "manga-name")

					// https://go.dev/blog/slices-intro
					items := make([]pterm.BulletListItem, len(mangas))

					for i, v := range mangas {
						title := fmt.Sprintf("(%d) - %s", i+1, v.Find("a").Text())
						items[i] = pterm.BulletListItem{Level: 0, Text: title,
							TextStyle: pterm.NewStyle(pterm.FgYellow)}
					}

					spinner.Stop()

					pterm.DefaultSection.Print("Results")
					pterm.DefaultBulletList.WithItems(items).Render()

					// Finds the url of the selected manga.
					val, _ := pterm.DefaultInteractiveTextInput.Show("Select index: ")
					index, _ := strconv.Atoi(val)

					url := fmt.Sprintf("https://mangareader.to%s", mangas[index-1].Find("a").Attrs()["href"])

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
