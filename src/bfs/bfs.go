package bfs

import (
	"fmt"
	"os"
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

	// Initiate array to store visited links, also acts as queue
	list := []string{}
	var visit_id uint = 0 // tracks id of currently visited link

	// Var to store current link
	var current_link string

	// Initiate bool to determine if solution is reached
	found := false

	// Initiate tree and id for tree
	t := tree.Empty[string]()
	var parent_id uint
	var current_id uint = 1 // tracks id for tree
	var final_id uint

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`^https://en.wikipedia.org/wiki/([^:]+)[^:]*$`),
		),
	)

	file, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("Failed to create log")
	}

	wikipediaRegex := regexp.MustCompile(`^/wiki/([^:]+)[^:]*$`)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absolute_link := e.Request.AbsoluteURL(link)

		if !found && (wikipediaRegex.MatchString(link)) && !(slices.Contains(list, absolute_link)) { // is a wikipedia link && not visited && not in queue
			// Print link
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			_, err = file.WriteString("Link found: " + e.Text + " -> " + link + "\n")

			// Append link to array of to-be visited links
			list = append(list, absolute_link)

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
		_, err = file.WriteString("Visited " + r.Request.URL.String() + "\n")
		link := r.Request.URL.String()
		if link == goalURL {
			found = true
			final_id = parent_id

			fmt.Println("Goal Found! " + link)
			_, err = file.WriteString("Goal Found! " + link + "\n")
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
	list = append(list, startURL)
	for (!found) && (int(visit_id) < len(list)) {
		current_link = list[visit_id]
		parent_id = visit_id
		visit_id++
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
