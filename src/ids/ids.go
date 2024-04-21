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

func IDS(startURL, goalURL string, currentDepth int, maxDepth int, visited *[]string, cek *int) {

	if maxDepth == currentDepth {
		*visited = append(*visited, startURL)

		return
	}

	if startURL == goalURL {
		*visited = append(*visited, startURL)
		*cek = 1
		return
	}


	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`^https://en.wikipedia.org/wiki/([^:]+)[^:]*$`),
		),
	)
	
	var tempArrayURL []string

	wikipediaRegex := regexp.MustCompile(`^/wiki/([^:]+)[^:]*$`)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if wikipediaRegex.MatchString(link) {
			tempArrayURL = append(tempArrayURL, e.Request.AbsoluteURL(link))
		}
	})

	err := c.Visit(startURL)

	if err != nil {
		fmt.Println("Can't visitiing: ", startURL)
		fmt.Println("Error: ", err)
	}

	
	*visited = append(*visited, startURL)


	for _, element := range tempArrayURL {

		if !contains(*visited, element) {

			IDS(element, goalURL, currentDepth, maxDepth-1, visited, cek)

			if (*visited)[len(*visited)-1] == goalURL {
				*cek = 1
				break
				} else if (*visited)[len(*visited)-1] != goalURL {
					*visited = (*visited)[:len(*visited)-1]
			}
		}
	}

}


func Main(startURL, goalURL string) *graph.Graph[string, string] {
	visited := make([]string, 0)

	var maxDepth int = 0

	var cek int = 0

	for cek == 0{

		visited := make([]string, 0)

		IDS(startURL, goalURL, 0, maxDepth, &visited, &cek)

		maxDepth += 1


		if cek == 1{

			break
		}
	}

	g := graph.New(graph.StringHash, graph.Directed())

	// fmt.Println("*")
	// fmt.Println("**")
	// for i := 0; i < len(visited);i++{
	// 	fmt.Println((visited)[i])
	// }
	// fmt.Println("**")
	// fmt.Println("*")


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
