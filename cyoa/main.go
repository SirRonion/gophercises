package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Story map[string]StoryArc

func getStory() (Story, error) {
	stories := make(map[string]StoryArc)
	data, err := os.ReadFile("gopher.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &stories); err != nil {
		return nil, err
	}

	return stories, nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		stories, err := getStory()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting data: %v", err), http.StatusInternalServerError)
			return
		}

		arcName := r.URL.Path[1:]
		if r.URL.Path == "/" {
			arcName = "intro"
		}

		data, exists := stories[arcName]
		if !exists {
			http.Error(w, "Error getting data: couldn't find story arc", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("./view/index.html")
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", mux)
}
