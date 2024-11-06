package main

import (
	// importing Colly
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {

	// instantiate a new collector object
	c := colly.NewCollector(
		colly.AllowedDomains("https://bikeindex.org/bikes?serial=&button=&location=you&distance=5&stolenness=proximity"),
	)

	// Go to this page, identify that location is still being mesasured (expect 154 results), divide by 10, if remainder, its ans+1 pages
	// Scrape 10 pages at the same time, write to a file, then move on to the next 10 pages (use mutexes and whatnot)

	// On every a element which has href attribute call callback
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
}
