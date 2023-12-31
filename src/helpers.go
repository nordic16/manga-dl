package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/pterm/pterm"
)

/* Searches for manga @query on mangareader.to  */
func scrapeMangas(query string) []soup.Root {
	resp, err := soup.Get("https://mangapill.com/search?q=" + query)

	if err != nil {
		os.Exit(1)
	}

	// Actually handles the requests.
	doc := soup.HTMLParse(resp)
	return doc.FindAll("a", "class", "mb-2")
}

/* Uses pterm to beautifully list mangas. */
func listMangas(mangas []soup.Root) {
	// https://go.dev/blog/slices-intro
	items := make([]pterm.BulletListItem, len(mangas))

	for i, v := range mangas {
		title := v.Find("div").Text()
		prompt := fmt.Sprintf("(%d) - %s", i+1, title)
		items[i] = pterm.BulletListItem{Level: 0, Text: prompt,
			BulletStyle: pterm.NewStyle(pterm.FgLightBlue)}
	}

	if len(items) > 0 {
		pterm.DefaultSection.Printf("Found %d Results", len(items))
		pterm.DefaultBulletList.WithItems(items).Render()
	}
}

/*
* Returns an array containing all chapters from @manga_url.
 */
func scrapeChapters(manga_url string) []string {
	resp, err := soup.Get(manga_url)
	chapters := []string{}

	if err != nil {
		os.Exit(-1)
	}

	doc := soup.HTMLParse(resp)
	chapterElems := doc.FindAll("a", "class", "border-border")

	for _, v := range chapterElems {
		name := v.Attrs()["href"]

		if !strings.Contains(name, ".") {
			chapters = append(chapters, name)
		}
	}

	return chapters

}

/* Scrapes all issues from a given chapter @chapter_url  */
func scrapeImages(chapter_url string) []string {
	req, e := soup.Get(chapter_url)

	if e != nil {
		pterm.Error.Print("Something went wrong!\n")
		os.Exit(-1)
	}

	doc := soup.HTMLParse(req)
	images := doc.FindAll("img", "class", "js-page")
	urls := make([]string, len(images))

	for i, v := range images {
		urls[i] = v.Attrs()["data-src"]
	}
	return urls
}

/* Downloads images */
func downloadImages(images []string) []string {
	pathToImages := make([]string, len(images))
	client := &http.Client{}
	defer client.CloseIdleConnections()
	dir, _ := os.MkdirTemp("", "manga-tmp")

	for i, imageUrl := range images {
		req, _ := http.NewRequest("GET", imageUrl, nil)
		req.Header.Set("Referer", "https://mangapill.com")

		img, e := client.Do(req)
		buf, _ := io.ReadAll(img.Body) // The actual raw data

		if e != nil {
			pterm.Error.Println("Something went wrong downloading the images!")
			os.Exit(-1)
		}

		name := ""

		// to get it to work with cbz.
		if i < 9 {
			name = fmt.Sprintf("img-0%d.jpeg", i+1)

		} else {
			name = fmt.Sprintf("img-%d.jpeg", i+1)
		}

		f, _ := os.Create(path.Join(dir, name))

		f.Write(buf)
		pathToImages[i] = f.Name()

		if e != nil {
			fmt.Println(e)
			os.Exit(-1)
		}
	}

	return pathToImages
}

func createCbzFile(dir string) string {
	fPath := "/tmp/comic.cbz"

	_, e := exec.LookPath("7z")

	if e != nil {
		pterm.Error.Println("Looks like 7z isn't installed... Install it and try again!")
	}

	cmd := fmt.Sprintf("7z a %s %s/*", fPath, dir)
	exec.Command("/bin/sh", "-c", cmd).Run()

	return fPath
}

// runes
func split(r rune) bool {
	return r == '-' || r == '.'
}
