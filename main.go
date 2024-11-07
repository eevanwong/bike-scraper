package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"
)

// This code snippet uses Playwright to scrape the Bike Index website.
// Playwright is a browser automation library that can be used to scrape websites. Playwright-go is just a Go wrapper around the Playwright library.
// I used this out of whatever was available and what I first saw, but in future, to make things more robust and go-oriented, use Rod
type Bike struct {
	Title      string
	Serial     string
	Colors     string
	DateStolen string
	Location   string
}

func main() {
	link := "https://bikeindex.org/bikes?serial=&button=&location=you&distance=5&stolenness=proximity"
	bikes := []Bike{}
	links := []string{}
	count := 0
	num_pages := 0

	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	// Launch the browser in non-headless mode
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		// Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	// Visit the page
	_, err = page.Goto(link)
	if err != nil {
		log.Fatalf("could not go to page: %v", err)
	}

	page.Locator(".close").First().Click()

	// Extract the count from the span with class "count"
	count_str, err := page.Locator("#stolenness_tab_proximity .count").First().TextContent()
	if err != nil {
		log.Fatalf("could not get text content: %v", err)
	}

	count_str = strings.Trim(count_str, "()")
	count, err = strconv.Atoi(count_str)

	num_pages = int(math.Ceil(float64(count) / 10))
	for i := 0; i < num_pages; i++ {
		link := link + "&page=" + strconv.Itoa(i+1)
		links = append(links, link)
		// fmt.Println(link)
	}

	for _, link := range links {
		_, err = page.Goto(link)
		if err != nil {
			log.Fatalf("could not go to page: %v", err)
		}
		new_bikes, err := scrapeBikes(page)
		if err != nil {
			log.Fatalf("could not scrape bikes: %v", err)
		}
		bikes = append(bikes, new_bikes...)
	}
	// Wait for user input to keep the browser open for debug
	// fmt.Println("Press Enter to close the browser...")
	// fmt.Scanln()
	// Close the browser
	if err := browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err := pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}

	// Write the bikes to a CSV file
	if err := writeBikesToCSV(bikes, "bikes.csv"); err != nil {
		log.Fatalf("could not write bikes to CSV: %v", err)
	}
}

func writeBikesToCSV(bikes []Bike, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"Title", "Serial", "Colors", "Date", "Location"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("could not write header: %v", err)
	}

	// Write the bike data
	for _, bike := range bikes {
		record := []string{bike.Title, bike.Serial, bike.Colors, bike.DateStolen, bike.Location}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("could not write record: %v", err)
		}
	}

	return nil
}

// detect number of "bike-box-item" per page
// divide find total no of pages

// grab h5 element with class ="title-link" -> get text
// there are 2 attr-lists per box, extract all of them
// each attr-list has 2 li items, the first li item is the title, with the second being the value (the classes and tags are diff so hard to differentiate)
func scrapeBikes(page playwright.Page) ([]Bike, error) {
	bikes := []Bike{}

	// Get the list of bikes
	bikeBoxes, err := page.Locator(".bike-box-item").All()
	if err != nil {
		return nil, fmt.Errorf("could not find bike boxes: %v", err)
	}

	// Iterate over each bike box
	for _, bikeBox := range bikeBoxes {
		bike := Bike{}

		// Get the title of the bike
		title, err := bikeBox.Locator(".title-link").First().TextContent()
		if err != nil {
			return nil, fmt.Errorf("could not find title: %v", err)
		}
		title = strings.TrimSpace(title)
		title = strings.Replace(title, "\n\n", " ", -1)

		fmt.Println(title)
		bike.Title = title

		// grab all li items
		lists, err := bikeBox.Locator("ul.attr-list").All()
		if err != nil {
			return nil, fmt.Errorf("could not find title: %v", err)
		}

		// 2 lists per box
		// fmt.Println(len(lists))

		// each list has 2 list items (1st is title, 2nd is value)
		// params := make(map[string]string)
		for _, list := range lists {
			// Get the title of the bike
			li, err := list.Locator("li").All()
			if err != nil {
				return nil, fmt.Errorf("could not find title: %v", err)
			}

			// each list has 2 list items (1st is title, 2nd is value)
			for _, item := range li {
				text, err := item.TextContent()
				if err != nil {
					return nil, fmt.Errorf("could not find title: %v", err)
				}

				text = strings.TrimSpace(text)
				text = strings.ReplaceAll(text, "\n", " ")
				text = strings.ToTitle(text)

				if strings.HasPrefix(text, "SERIAL") {
					text = strings.Replace(text, "SERIAL: ", "", 1)
					bike.Serial = text
				} else if strings.HasPrefix(text, "PRIMARY COLORS") {
					text = strings.Replace(text, "PRIMARY COLORS: ", "", 1)
					bike.Colors = text
				} else if strings.HasPrefix(text, "STOLEN") {
					text = strings.Replace(text, "STOLEN: ", "", 1)
					bike.DateStolen = text
				} else if strings.HasPrefix(text, "LOCATION") {
					text = strings.Replace(text, "LOCATION: ", "", 1)
					bike.Location = text
				}
			}
		}

		bikes = append(bikes, bike)
	}

	return bikes, nil
}
