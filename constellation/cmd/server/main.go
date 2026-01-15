package main

import (
	"fmt"
	"log"
	"net/http"

	"constellation"
)

var c *constellation.Constellation

func main() {
	// Generate constellation once at startup
	c = constellation.GenerateConstellation(100, 15, 1000, 500)
	fmt.Printf("Generated constellation: %d nodes, %d edges\n", len(c.Storage), countEdges(c))

	// Serve static files (viewer.html)
	http.Handle("/", http.FileServer(http.Dir(".")))

	// API endpoint for routes
	http.HandleFunc("/api/data", handleData)

	fmt.Println("Server running at http://localhost:9000")
	fmt.Println("Try: http://localhost:9000/viewer.html?from=GS-0&to=GS-5")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handleData(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	var data []byte
	var err error

	if from != "" && to != "" {
		data, err = c.ExportJSONWithRoute(from, to)
	} else {
		data, err = c.ExportJSON()
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func countEdges(c *constellation.Constellation) int {
	count := 0
	for _, neighbors := range c.AdjacentList {
		count += len(neighbors)
	}
	return count / 2
}
