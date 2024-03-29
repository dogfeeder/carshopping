package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	fName := "carresults.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Year", "Title", "Odometer", "Price", "State"})

	// Instantiate default collector
	c := colly.NewCollector(
		colly.UserAgent(""),
	)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.90 Safari/537.36")
		r.Headers.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;")
		r.Headers.Set("Accept-Encoding", "gzip")
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response", r.StatusCode)
	})

	c.OnHTML(".listing-item", func(e *colly.HTMLElement) {
		e.DOM.Find("div.n_width-max.title > a > h2").Children().Remove()
		fullTitle := e.ChildText("div.n_width-max.title > a > h2")
		year := fullTitle[0:4]
		title := fullTitle[5:]
		odometer := e.ChildText("div:nth-child(1) > div.feature-text")
		price := strings.TrimRight(e.ChildText(".price"), "*")
		state := strings.TrimRight(e.ChildText("div.franchise-name"), "-")
		state = strings.TrimSpace(state)
		writer.Write([]string{
			year,
			title,
			odometer,
			price,
			state,
		})
		log.Printf("Car Found: %s - %s, %s, %s", title, odometer, price, state)
	})

	c.OnHTML(".next", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		e.Request.Visit(link)
	})

	// c.Visit("https://www.carsales.com.au/cars/results/?q=%28And.Service.Carsales._.%28C.Make.BMW._.%28C.MarketingGroup.1%20Series._.Model.135i.%29%29_.GenericGearType.Manual._.BodyStyle.Coupe._.%28Or.State.Queensland._.State.New%20South%20Wales.%29%29&sortby=~Price&WT.z_srchsrcx=makemodel")
	c.Visit("https://www.carsales.com.au/cars/results/?q=%28And.Service.Carsales._.GenericGearType.Manual._.BodyStyle.Coupe._.%28Or.State.Queensland._.State.New%20South%20Wales.%29%29&sortby=~Price&WT.z_srchsrcx=makemodel")
	log.Printf("Scraping finished, check file %q for results\n", fName)
}
