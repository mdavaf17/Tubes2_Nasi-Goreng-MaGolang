package bfs

import (
	"fmt"
	"os"
	"regexp"
	"slices"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/gocolly/colly"
)

func Main(startURL, goalURL string) *graph.Graph[string, string] {
	fmt.Println("BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Goal URL:", goalURL)

	// Initiate array to store visited links
	visited := []string{}

	// Initiate queue to store to-be visited links
	queue := []string{}

	// Var to store current link
	var current_link string

	// Initiate bool to determine if solution is reached
	found := false

	// Initiate graph
	g := graph.New(graph.StringHash, graph.Directed())

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`^https://en.wikipedia.org/wiki/([^:]+)[^:]*$`),
		),
	)

	wikipediaRegex := regexp.MustCompile(`^/wiki/([^:]+)[^:]*$`)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if (wikipediaRegex.MatchString(link)) && !(slices.Contains(visited, link)) && !(slices.Contains(queue, link)) { // is a wikipedia link && not visited && not in queue
			// Print link
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)

			// Append link to array of to-be visited links
			absolute_link := e.Request.AbsoluteURL(link)
			queue = append(queue, absolute_link)

			_ = g.AddVertex(absolute_link)

			_ = g.AddEdge(current_link, absolute_link)

		}
		// Visit link found on page
		// Only those links are visited which are matched by  any of the URLFilter regexps
		// if len(queue) > 0 {
		// 	visit := queue[0]
		// 	queue = queue[1:]
		// 	c.Visit(e.Request.AbsoluteURL(visit))
		// }
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on start_url
	queue = append(queue, startURL)
	_ = g.AddVertex(startURL)
	for (!found) && (len(queue) > 0) {
		current_link = queue[0]
		if current_link == goalURL {
			found = true
			fmt.Println("Link Found! " + current_link)

		} else {
			queue = queue[1:]
			visited = append(visited, current_link)
			c.Visit(current_link)
		}
	}

	fmt.Println(startURL)
	file, _ := os.Create("./mygraph.gv")
	_ = draw.DOT(g, file)

	// find path from graph
	path_all, _ := graph.AllPathsBetween(g, startURL, goalURL)
	path, _ := graph.ShortestPath(g, startURL, goalURL)
	fmt.Println(path_all)
	fmt.Println(path)

	return &g
}
