package ids

import (
	"fmt"
	"regexp"

	"github.com/dominikbraun/graph"
	"github.com/gocolly/colly"
)

func contains(arr []string, target string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

func IDS(startURL, goalURL string, currentDepth int, maxDepth int, visited *[]string) {
	if maxDepth == currentDepth {
		*visited = append(*visited, startURL)
		if startURL == goalURL {
			return
		} else {
			(*visited) = (*visited)[:len(*visited)-1]
		}
		return
	}

	if startURL == goalURL {
		*visited = append(*visited, startURL)
		return
	}

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`^https://en.wikipedia.org/wiki/([^:]+)[^:]*$`),
		),
	)

	err := c.Visit(startURL)

	if err != nil {
		fmt.Println("Can't visitiing: ", startURL)
		fmt.Println("Error: ", err)
	}

	*visited = append(*visited, startURL)

	var tempArrayURL []string

	wikipediaRegex := regexp.MustCompile(`^/wiki/([^:]+)[^:]*$`)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if wikipediaRegex.MatchString(link) {
			tempArrayURL = append(tempArrayURL, e.Request.AbsoluteURL(link))
		}
	})

	for i := 0; i < len(tempArrayURL); i++ {
		if !contains(*visited, tempArrayURL[i]) {
			IDS(tempArrayURL[i], goalURL, currentDepth, maxDepth-1, visited)

			if (*visited)[len(*visited)-1] == goalURL {
				break
			} else if (*visited)[len(*visited)-1] != goalURL {
				*visited = (*visited)[:len(*visited)-1]
			}
		}
	}
}

func Main(startURL, goalURL string) *graph.Graph[string, string] {
	visited := make([]string, 0)

	maxDepth := 0

	for len(visited) == 0 || visited[len(visited)-1] != goalURL {
		IDS(startURL, goalURL, 0, maxDepth, &visited)
		maxDepth += 1
		if(visited[len(visited)-1] == goalURL){
			break
		}
	}

	g := graph.New(graph.StringHash, graph.Directed())

	for i := 0; i < len(visited); i++ {
		_ = g.AddVertex(visited[i])
	}

	for i := 0; i < len(visited); i++ {
		if i+1 < len(visited) {
			_ = g.AddEdge(visited[i], visited[i+1])
		}
	}

	return &g
}
