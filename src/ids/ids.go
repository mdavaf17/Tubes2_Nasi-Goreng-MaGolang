package ids

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/gocolly/colly"
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

func contains(arr []string, target string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
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

func IDS(startURL, goalURL string, currentDepth int, maxDepth int, visited *[]string, cek *int, periksa *map[string]bool, visited_heuristic *[]string) {

	if maxDepth == 0 {
		if len(*visited) == 0 {
			if goalURL == startURL {
				*visited = append(*visited, startURL)
				return
			}
		}
	}

	if currentDepth >= maxDepth {
		*visited = append(*visited, startURL)

		return
	}

	if startURL == goalURL {
		*visited = append(*visited, startURL)
		*cek = 1
		return
	}

	c := colly.NewCollector()

	var tempArrayURL []string

	wikipediaRegex := regexp.MustCompile(`^https://en.wikipedia.org/wiki/`)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link, _, _ := strings.Cut(e.Attr("href"), "#")
		absolute_link := e.Request.AbsoluteURL(link)
		if (wikipediaRegex.MatchString(absolute_link)) && (absolute_link != "https://en.wikipedia.org/wiki/Main_Page") && doesNotContainAnyPrefix(absolute_link) {
			tempArrayURL = append(tempArrayURL, e.Request.AbsoluteURL(link))
		}
	})

	// var wg sync.WaitGroup
	// var mu sync.Mutex

	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {

	// 	link := e.Attr("href")
	// 	if wikipediaRegex.MatchString(link) {
	// 		wg.Add(1)
	// 		go func(link string) {
	// 				defer wg.Done()
	// 				mu.Lock()
	// 				defer mu.Unlock()
	// 				tempArrayURL = append(tempArrayURL, e.Request.AbsoluteURL(link))
	// 		}(link)

	// 	}
	// })

	// wg.Wait()

	err := c.Visit(startURL)
	(*periksa)[startURL] = true

	if err != nil {
		fmt.Println("Can't visitiing: ", startURL)
		fmt.Println("Error: ", err)
	}

	*visited = append(*visited, startURL)

	*visited_heuristic = append(*visited_heuristic, startURL)

	for _, element := range tempArrayURL {
		// fmt.Println(element)
		// fmt.Println(goalURL)

		if !contains(*visited_heuristic, element) {

			IDS(element, goalURL, currentDepth+1, maxDepth, visited, cek, periksa, visited_heuristic)

			if (*visited)[len(*visited)-1] == goalURL {
				*cek = 1
				break
			} else if (*visited)[len(*visited)-1] != goalURL {
				*visited = (*visited)[:len(*visited)-1]
			}
		}
	}

}

func Main(startURL, goalURL string) (*graph.Graph[string, string], int) {
	visited := []string{}

	periksa := make(map[string]bool)

	var maxDepth int = 0

	var cek int = 0

	// fmt.Println(startURL)
	// fmt.Println(goalURL)
	visited_heuristic := []string{}

	for cek == 0 {

		visited = []string{}

		visited_heuristic = []string{}

		IDS(startURL, goalURL, 0, maxDepth, &visited, &cek, &periksa, &visited_heuristic)

		maxDepth += 1

		if cek == 1 {

			break
		}
	}

	g := graph.New(graph.StringHash, graph.Directed())

	fmt.Println("*")
	fmt.Println("**")
	for i := 0; i < len(visited); i++ {
		fmt.Println((visited)[i])
	}
	fmt.Println("**")
	fmt.Println("*")

	for _, v := range visited {
		_ = g.AddVertex(getWikipediaTitleFromURL(v), graph.VertexAttribute("URL", v))
	}

	// Add edges to the graph
	for i := 0; i+1 < len(visited); i++ {
		_ = g.AddEdge(getWikipediaTitleFromURL(visited[i]), getWikipediaTitleFromURL(visited[i+1]))
	}

	return &g, len(periksa)
}
