package bfs

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/gocolly/colly"
	"github.com/kingledion/go-tools/tree"
)

func Main(startURL, goalURL string) (*graph.Graph[string, string], int) {
	// file, err := os.Create("log.txt")
	// if err != nil {
	// 	fmt.Println("Failed to create log")
	// }

	fmt.Println("BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Goal URL:", goalURL)
	// _, err = file.WriteString("Start URL: " + startURL)
	// _, err = file.WriteString("Goal URL: " + goalURL)

	// Initiate array to store visited links, also acts as queue
	list := []string{}
	var visit_id uint = 0 // tracks id of currently visited link

	// Array to store banned links
	black_list := []string{"https://en.wikipedia.org/wiki/Main_Page"}

	// Initiate bool to determine if solution is reached
	found := false

	// Initiate tree and id for tree
	t := tree.Empty[string]()
	var parent_id uint
	var current_id uint = 1 // tracks id for tree
	var final_id uint

	// Counter for visited links
	num_visited := 0

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`^https://en.wikipedia.org/wiki/([^:]+)[^:]*$`),
		),
	)

	wikipediaRegex := regexp.MustCompile(`^/wiki/([^:]+)[^:]*$`)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// link := e.Attr("href")
		link, _, _ := strings.Cut(e.Attr("href"), "#")
		absolute_link := e.Request.AbsoluteURL(link)

		if !found && (wikipediaRegex.MatchString(link)) && !(slices.Contains(list, absolute_link)) && !(slices.Contains(black_list, absolute_link)) { // is a wikipedia link && not visited && not in queue

			// Print link
			// fmt.Printf("Link found: %q -> %s : %s\n", e.Text, link)
			// _, err = file.WriteString("Link found: " + e.Text + " -> " + link + "\n")

			// Append link to array of to-be visited links
			list = append(list, absolute_link)

			t.Add(current_id, parent_id, absolute_link)

			if absolute_link == goalURL {
				found = true
				fmt.Println("Goal Found! " + link)
				// _, err = file.WriteString("Goal Found! " + link + "\n")
				final_id = current_id
			}

			current_id++
		}
	})

	// Before checking HTML check if link is goal
	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL.String()
		articleID := "123"

		fmt.Println("Visited " + r.Request.URL.String() + " : " + articleID)
		// _, err = file.WriteString("Visited " + r.Request.URL.String() + " : " + articleID + "\n")

		if link == goalURL {
			found = true
			final_id = parent_id

			fmt.Println("Goal Found! " + link)
			// _, err = file.WriteString("Goal Found! " + link + "\n")
		}

		num_visited++
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

	// Var to store current link
	var current_link string

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

	// Output jumlah artikel yang diperiksa dan dilalui
	fmt.Println("Jumlah artikel yang diperiksa:", num_visited)
	fmt.Println("Jumlah artikel yang dilalui:", len(path))
	// _, err = file.WriteString("Jumlah artikel yang diperiksa: " + string(num_visited))
	// _, err = file.WriteString("Jumlah artikel yang dilalui: " + string(len(path)))

	return &g, num_visited
}
