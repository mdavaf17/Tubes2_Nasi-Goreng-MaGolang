package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/go-zoox/fetch"
	"github.com/gocolly/colly"
	// "github.com/mdavaf17/Tubes2_Nasi-Goreng-MaGolang/src/tree"
)

type Article struct {
	Title string
	URL   string
}

func main() {
	fmt.Println("Apllication is running")

	home := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))

		tmpl.Execute(w, nil)
	}

	search := func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("id")
		var search_title string

		if ID == "inputStart" {
			search_title = r.PostFormValue("inputStartTitle")
		} else {
			search_title = r.PostFormValue("inputGoalTitle")
		}

		if search_title != "" {
			url := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=opensearch&search=%s&limit=10&namespace=0&format=json", search_title)
			response, _ := fetch.Get(url)
			jsonData, _ := response.JSON()

			// Unmarshal JSON data
			var data []interface{}
			if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
				panic(err)
			}

			// Extract titles and URLs from the parsed data
			titles := data[1].([]interface{})
			urls := data[3].([]interface{})

			// Construct Article structs
			var articles []Article
			for i := 0; i < len(titles); i++ {
				title := titles[i].(string)
				url := urls[i].(string)
				article := Article{
					Title: title,
					URL:   url,
				}
				articles = append(articles, article)
			}

			// Define the template for list items
			tmpl := template.Must(template.New("listItems").Parse(`
				{{$ID := .ID}}
				{{range .Articles}}
					<div class="suggestion-item" onclick="setSearchInput('{{$ID}}', '{{.Title}}', '{{.URL}}')">
						{{.Title}} <sub>{{.URL}}</sub>
					</div>
				{{end}}
			`))

			// Create a struct to hold both ID and Articles
			type TemplateData struct {
				ID       string
				Articles []Article
			}

			// Populate TemplateData struct
			tmplData := TemplateData{
				ID:       ID, // Pass the ID to the template
				Articles: articles,
			}

			// Execute the template with tmplData
			tmpl.Execute(w, tmplData)
		}
	}

	race := func(w http.ResponseWriter, r *http.Request) {
		start_title := r.PostFormValue("inputStartTitle")
		start_url := r.PostFormValue("inputStartURL")
		goal_title := r.PostFormValue("inputGoalTitle")
		goal_url := r.PostFormValue("inputGoalURL")

		fmt.Println(start_title, start_url, goal_title, goal_url)

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

				// Visit link found on page
				// Only those links are visited which are matched by  any of the URLFilter regexps
				c.Visit(e.Request.AbsoluteURL(link))
			}
		})

		// Before making a request print "Visiting ..."
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		// Start scraping on http://httpbin.org
		c.Visit(start_url)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/search/", search)
	http.HandleFunc("/race/", race)

	// Exit when error
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
