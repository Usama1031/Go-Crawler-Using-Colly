package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []movie
}

type movie struct {
	Title string
	Year  string
}

func main() {
	month := flag.Int("month", 1, "Month to fetch birthdays for")
	day := flag.Int("day", 1, "Day to fetch birthdays for")
	flag.Parse()
	crawl(*month, *day)

}

func crawl(month int, day int) {

	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	infoCollector := c.Clone()

	c.OnHTML(".ipc-avatar", func(e *colly.HTMLElement) {
		profileUrl := e.ChildAttr("a.ipc-lockup-overlay.ipc-focusable", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)

		infoCollector.Visit(profileUrl)
	})

	// c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
	// 	nextPage := e.Request.AbsoluteURL(e.Attr("href"))
	// 	c.Visit(nextPage)
	// })

	infoCollector.OnHTML(".ipc-page-content-container", func(e *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = e.ChildText("h1.sc-ec65ba05-0 span.hero__primary-text")
		tmpProfile.Photo = e.ChildAttr("div.sc-9a2a0028-7  div.ipc-media img.ipc-image", "src")
		tmpProfile.JobTitle = e.ChildText("div.sc-78c11d06-0 ul.ipc-inline-list li.ipc-inline-list__item")
		tmpProfile.BirthDate = e.ChildText("div.sc-59a43f1c-1 span.sc-59a43f1c-2:nth-of-type(2)")
		tmpProfile.Bio = strings.TrimSpace(e.ChildText("div.ipc-html-content div.ipc-html-content-inner-div"))

		e.ForEach("div.ipc-list-card--span", func(_ int, kf *colly.HTMLElement) {
			tmpMovie := movie{}
			tmpMovie.Title = kf.ChildText("div.ipc-primary-image-list-card__content div.ipc-primary-image-list-card__content-top a.ipc-primary-image-list-card__title")
			tmpMovie.Year = kf.ChildText("div.ipc-primary-image-list-card__content div.ipc-primary-image-list-card__content-bottom span.ipc-primary-image-list-card__secondary-text")
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		})

		js, err := json.MarshalIndent(tmpProfile, "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitng:", r.URL.String())
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visting profile URL", r.URL.String())
	})

	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	c.Visit(startUrl)
}
