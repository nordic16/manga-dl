package main

import (
	"fmt"
	"os"
	"strconv"

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

/*
	 Scrapes the number of chapters. Due to the nature of mangareader.to, only the number of chapters is required
		in order to form an url for a given chapter.
		Structure: mangareader.to/read/<manga-id>/en/chapter-<num>
*/
func totalChapters(manga_url string) int {
	resp, err := soup.Get(manga_url)

	if err != nil {
		os.Exit(-1)
	}

	doc := soup.HTMLParse(resp)
	latestChapterNum := doc.Find("ul", "id", "en-chapters").Find("li", "class", "reading-item").Attrs()["data-number"]
	num, _ := strconv.Atoi(latestChapterNum)

	return num

}
