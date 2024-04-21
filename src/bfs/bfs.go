package bfs

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/dominikbraun/graph"
	"github.com/gocolly/colly"
	"github.com/kingledion/go-tools/tree"
)

func Main(startURL, goalURL string) *graph.Graph[string, string] {
	fmt.Println("BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Goal URL:", goalURL)

	// Initiate array to store visited links
	visited := []string{}

	// Initiate queue to store to-be visited links
	queue := []string{}
	queue_id := []uint{}

	// Var to store current link
	var current_link string

	// Initiate bool to determine if solution is reached
	found := false

	// Initiate tree and id for tree
	t := tree.Empty[string]()
	var parent_id uint
	var current_id uint = 1

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
			queue_id = append(queue_id, current_id)

			t.Add(current_id, parent_id, absolute_link)

			current_id++
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Find path to goal URL
	findPath := func(id uint) []string {
		path := []string{}

		n, _ := t.Find(id)
		for !(n.GetID() == 0) {
			path = append(path, n.GetData())
			n = n.GetParent()
		}

		path = append(path, n.GetData())

		return path
	}

	// Start scraping on start_url
	t.Add(0, 0, startURL)
	queue = append(queue, startURL)
	queue_id = append(queue_id, 0)
	for (!found) && (len(queue) > 0) {
		current_link = queue[0]
		parent_id = queue_id[0]
		if current_link == goalURL {
			found = true
			fmt.Println("Link Found! " + current_link)

		} else {
			queue = queue[1:]
			queue_id = queue_id[1:]
			visited = append(visited, current_link)
			c.Visit(current_link)
		}
	}

	// Get path to goal URL
	path := findPath(parent_id)

	// Initiate graph
	g := graph.New(graph.StringHash, graph.Directed())

	// Add path to graph
	for i := len(path) - 1; i > 0; i-- {
		_ = g.AddVertex(path[i])
		_ = g.AddVertex(path[i-1])
		_ = g.AddEdge(path[i], path[i-1])
	}

	return &g
}
