package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/go-zoox/fetch"
	"github.com/mdavaf17/Tubes2_Nasi-Goreng-MaGolang/src/bfs"
	"github.com/mdavaf17/Tubes2_Nasi-Goreng-MaGolang/src/ids"
)

type Article struct {
	Title string
	URL   string
}

type TemplateData struct {
	ID       string
	Articles []Article
}

func main() {
	fmt.Println("Apllication is running in port 8030")

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
			url := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=opensearch&search=%s&limit=10&namespace=0&format=json&redirects=resolve", search_title)
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
		// start_title := r.PostFormValue("inputStartTitle")
		start_url := r.PostFormValue("inputStartURL")
		// goal_title := r.PostFormValue("inputGoalTitle")
		goal_url := r.PostFormValue("inputGoalURL")
		algorithm := r.PostFormValue("inputAlgorithm")

		var graphResult *graph.Graph[string, string]
		var numChecked int

		t1 := time.Now()
		if algorithm == "IDS" {
			graphResult, numChecked = ids.Main(start_url, goal_url)
		} else {
			graphResult, numChecked = bfs.Main(start_url, goal_url)
		}
		t2 := time.Now()
		duration := t2.Sub(t1)
		minutes := int(duration.Minutes())
		seconds := duration.Seconds() - float64(minutes)*60
		lenRes, _ := (*graphResult).Order()
		lenRes -= 1

		// GRAPH RESULT
		var buf bytes.Buffer
		draw.DOT(*graphResult, &buf)

		resDOT := buf.String()

		tmpl := template.Must(template.New("graphItems").Parse(`
		<script type="text/javascript">
			var dot = {{.ResDOT}};
			var algo = {{.Algo}};
			var checked = {{.LenChecked}};
			var vertices = {{.LenRes}};
			var options = {
				format: 'svg',
			};
		
			var image = Viz(dot, options);
			var outputElement = document.getElementById('output');
		
			var algorithmElement = document.createElement('h5');
			algorithmElement.textContent = algo + ' Output';
		
			var listElement = document.createElement('ul');
			listElement.className = "list-group";
		
			// Create list items and set their inner HTML with the computed values
			var listItem0 = document.createElement('li');
			listItem0.innerHTML = image;
		
			var listItem1 = document.createElement('li');
			listItem1.innerHTML = 'Number of checked article: ' + checked;
		
			var listItem2 = document.createElement('li');
			listItem2.innerHTML = 'Number of article in solution: ' + vertices;
		
			var listItem3 = document.createElement('li');
			listItem3.innerHTML = 'Time: {{.Minute}} minutes ' + parseFloat({{.Second}}).toFixed(5) + ' seconds';
		
			// Append the list items to the unordered list
			listItem0.className = "list-group-item justify-content-center d-flex";
			listItem1.className = "list-group-item";
			listItem2.className = "list-group-item";
			listItem3.className = "list-group-item";
			listElement.appendChild(listItem0);
			listElement.appendChild(listItem1);
			listElement.appendChild(listItem2);
			listElement.appendChild(listItem3);
		
			// Append elements to outputElement
			outputElement.appendChild(algorithmElement);
			outputElement.appendChild(listElement);
			stopTimer();
		</script>	
			`))

		tmplData := struct {
			ResDOT     string
			Algo       string
			LenChecked int
			LenRes     int
			Minute     int
			Second     float64
		}{
			ResDOT:     resDOT,
			Algo:       algorithm,
			LenChecked: numChecked,
			LenRes:     lenRes,
			Minute:     minutes,
			Second:     seconds,
		}

		tmpl.Execute(w, tmplData)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/search/", search)
	http.HandleFunc("/race/", race)

	// Exit when error
	log.Fatal(http.ListenAndServe(":8030", nil))
}
