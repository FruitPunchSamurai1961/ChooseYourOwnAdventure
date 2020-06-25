package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Story map[string]Chapter

func exit(err error) {
	fmt.Printf("failed to read the file %v", err)
	os.Exit(1)
}

func parseFile() Story {
	file, err := os.Open("gopher.json")

	if err != nil {
		exit(err)
	}
	var story Story
	d := json.NewDecoder(file)

	if wrong := d.Decode(&story); wrong != nil {
		exit(err)
	}
	return story
}

func main() {
	story := parseFile()
	r := mux.NewRouter()
	port := os.Getenv("PORT")
	r.HandleFunc(":"+port, func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, ":"+port+"adventure/intro", http.StatusFound)
	})
	tmpl := template.Must(template.ParseFiles("layout.html"))
	r.HandleFunc("/adventure/{chapter}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		chosenChapter := vars["chapter"]
		chapter := story[chosenChapter]
		tmpl.Execute(writer, chapter)
	})

	http.ListenAndServe(":"+port, r)
}
