package main

import (
	"fmt"
	"os"

	"github.com/anaskhan96/soup"
	"github.com/pterm/pterm"
)

/* Searches for manga @query on mangareader.to  */
func scrapeMangas(query string) []soup.Root {
	resp, err := soup.Get("https://mangareader.to/search?keyword=" + query)

	if err != nil {
		os.Exit(1)
	}

	// Actually handles the requests.
	doc := soup.HTMLParse(resp)
	return doc.FindAll("h3", "class", "manga-name")
}

/* Uses pterm to beautifully list mangas. */
func listMangas(mangas []soup.Root) {
	// https://go.dev/blog/slices-intro
	items := make([]pterm.BulletListItem, len(mangas))

	for i, v := range mangas {
		title := fmt.Sprintf("(%d) - %s", i+1, v.Find("a").Text())
		items[i] = pterm.BulletListItem{Level: 0, Text: title,
			TextStyle: pterm.NewStyle(pterm.FgYellow)}
	}

	if len(items) > 0 {
		pterm.DefaultSection.Printf("Found %d Results", len(items))
		pterm.DefaultBulletList.WithItems(items).Render()
	}
}

/* Returns url */
func getUrl(manga soup.Root) string {
	return manga.Find("a").Attrs()["href"]
}
