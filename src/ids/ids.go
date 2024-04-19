package ids

import (
	"fmt"
	"regexp"

	"github.com/dominikbraun/graph"
	"github.com/gocolly/colly"
)

func Main(startURL, goalURL string) *graph.Graph[string, string] {
	fmt.Println("IDS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Goal URL:", goalURL)

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`^https://en.wikipedia.org/wiki/([^:]+)[^:]*$`),
		),
	)

	wikipediaRegex := regexp.MustCompile(`^/wiki/([^:]+)[^:]*$`)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if wikipediaRegex.MatchString(link) {
			// Print link
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		}
		// Visit link found on page
		// Only those links are visited which are matched by  any of the URLFilter regexps
		// c.Visit(e.Request.AbsoluteURL(link))

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on start_url
	// c.Visit(start_url)
	// c.Visit("https://linktr.ee/RPL2024")

	g := graph.New(graph.StringHash, graph.Directed())

	_ = g.AddVertex("Polandia")
	_ = g.AddVertex("B")
	_ = g.AddVertex("C")

	_ = g.AddEdge("Polandia", "C")
	_ = g.AddEdge("Polandia", "B")

	return &g
}
