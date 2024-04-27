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

func doesNotContainAnyPrefix(s string) bool {
	prefixes := []string{
		"https://en.wikipedia.org/wiki/Talk:",
		"https://en.wikipedia.org/wiki/User:",
		"https://en.wikipedia.org/wiki/User talk:",
		"https://en.wikipedia.org/wiki/Wikipedia:",
		"https://en.wikipedia.org/wiki/WP:",
		"https://en.wikipedia.org/wiki/WT:",
		"https://en.wikipedia.org/wiki/Wikipedia talk:",
		"https://en.wikipedia.org/wiki/File:",
		"https://en.wikipedia.org/wiki/File talk:",
		"https://en.wikipedia.org/wiki/MediaWiki:",
		"https://en.wikipedia.org/wiki/MediaWiki talk:",
		"https://en.wikipedia.org/wiki/Template:",
		"https://en.wikipedia.org/wiki/Template talk:",
		"https://en.wikipedia.org/wiki/Help:",
		"https://en.wikipedia.org/wiki/Help talk:",
		"https://en.wikipedia.org/wiki/Category:",
		"https://en.wikipedia.org/wiki/Category:talk",
		"https://en.wikipedia.org/wiki/Portal:",
		"https://en.wikipedia.org/wiki/Portal talk:",
		"https://en.wikipedia.org/wiki/Draft:",
		"https://en.wikipedia.org/wiki/Draft talk:",
		"https://en.wikipedia.org/wiki/TimedText:",
		"https://en.wikipedia.org/wiki/TimedText talk:",
		"https://en.wikipedia.org/wiki/Module:",
		"https://en.wikipedia.org/wiki/Module talk:",
		"https://en.wikipedia.org/wiki/Image:",
		"https://en.wikipedia.org/wiki/Image Talk:",
		"https://en.wikipedia.org/wiki/Topic:",
		"https://en.wikipedia.org/wiki/Special:",
		"https://en.wikipedia.org/wiki/Media:",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return false
		}
	}
	return true
}

func getWikipediaTitleFromURL(url string) string {
	var title string

	// Instantiate a new collector
	co := colly.NewCollector()

	// Find the title element
	co.OnHTML("title", func(e *colly.HTMLElement) {
		title = strings.Split(e.Text, " - Wikipedia")[0]
	})

	// Visit the URL
	co.Visit(url)

	return title
}

func Main(startURL, goalURL string) (*graph.Graph[string, string], int) {
	// Initiate array to store visited links, also acts as queue
	list := []string{}
	var visit_id uint = 0 // tracks id of currently visited link

	// Initiate bool to determine if solution is reached
	found := false

	// Initiate tree and id for tree
	t := tree.Empty[string]()
	var parent_id uint
	var current_id uint = 1 // tracks id for tree
	var final_id uint

	// Counter for visited links
	num_visited := 0

	c := colly.NewCollector()

	wikipediaRegex := regexp.MustCompile(`^https://en.wikipedia.org/wiki/`)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if found {
			return
		}

		// link := e.Attr("href")
		link, _, _ := strings.Cut(e.Attr("href"), "#")
		absolute_link := e.Request.AbsoluteURL(link)

		if (wikipediaRegex.MatchString(absolute_link)) && (absolute_link != "https://en.wikipedia.org/wiki/Main_Page") && doesNotContainAnyPrefix(absolute_link) && !slices.Contains(list, absolute_link) { // is a wikipedia link && not visited && not in queue

			// Append link to array of to-be visited links
			list = append(list, absolute_link)

			t.Add(current_id, parent_id, absolute_link)

			if absolute_link == goalURL {
				found = true
				final_id = current_id
			}

			current_id++
		}
	})

	// Before checking HTML check if link is goal
	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL.String()

		if link == goalURL {
			found = true
			final_id = parent_id
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
		err := c.Visit(current_link)

		if err != nil {
			fmt.Println("Can't visit: ", current_link)
			fmt.Println("Error: ", err)
		}
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
	_ = g.AddVertex(getWikipediaTitleFromURL(path[len(path)-1]))
	for i := len(path) - 2; i > -1; i-- {
		_ = g.AddVertex(getWikipediaTitleFromURL(path[i]), graph.VertexAttribute("URL", path[i]))
		_ = g.AddEdge(getWikipediaTitleFromURL(path[i+1]), getWikipediaTitleFromURL(path[i]))
	}

	return &g, num_visited
}
