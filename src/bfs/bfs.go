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
	var final_id uint

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`^https://en.wikipedia.org/wiki/([^:]+)[^:]*$`),
		),
	)

	wikipediaRegex := regexp.MustCompile(`^/wiki/([^:]+)[^:]*$`)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absolute_link := e.Request.AbsoluteURL(link)

		if !found && (wikipediaRegex.MatchString(link)) && !(slices.Contains(visited, absolute_link)) && !(slices.Contains(queue, absolute_link)) { // is a wikipedia link && not visited && not in queue
			// Print link
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)

			// Append link to array of to-be visited links
			queue = append(queue, absolute_link)
			queue_id = append(queue_id, current_id)

			t.Add(current_id, parent_id, absolute_link)

			// if absolute_link == goalURL {
			// 	found = true
			// 	fmt.Printf("Link Found! " + absolute_link)
			// 	final_id = current_id
			// }

			current_id++
		}
	})

	// Before checking HTML check if link is goal
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL.String())
		link := r.Request.URL.String()
		if link == goalURL {
			found = true
			fmt.Println("Link Found! " + link)
			final_id = parent_id
		}
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
		queue = queue[1:]
		queue_id = queue_id[1:]
		visited = append(visited, current_link)
		c.Visit(current_link)
	}

	// Get path to goal URL
	path := []string{}
	if found {
		path = findPath(final_id)
		fmt.Println(path)
	}

	// Initiate graph
	g := graph.New(graph.StringHash, graph.Directed())

	// Add path to graph
	_ = g.AddVertex(path[len(path)-1])
	for i := len(path) - 2; i > -1; i-- {
		_ = g.AddVertex(path[i])
		_ = g.AddEdge(path[i+1], path[i])
	}

	return &g
}
