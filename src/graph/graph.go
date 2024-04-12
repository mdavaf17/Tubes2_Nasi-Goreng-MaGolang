package main

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

func main() {
	g := graph.New(graph.IntHash, graph.Directed())

	_ = g.AddVertex(1)
	_ = g.AddVertex(2)
	_ = g.AddVertex(3)

	_ = g.AddEdge(1, 2)
	_ = g.AddEdge(1, 3)

	var buf bytes.Buffer
	err := draw.DOT(g, &buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	dotString := buf.String()

	// Remove line breaks and extra spaces
	re := regexp.MustCompile(`\s+`)
	dotString = re.ReplaceAllString(dotString, " ")

	fmt.Println(dotString)
}
