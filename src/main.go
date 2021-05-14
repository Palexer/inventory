package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

// rootTemplate references the specified rootTemplate and caches the parsed results
// to help speed up response times.
var rootTemplate = template.Must(template.ParseFiles("./templates/base.html", "./templates/index.html"))

var data = csvData{
	contentPath: "inventory_data.csv",
	cachePath:   "inventory_cache.csv",
}

func public() http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// load data
	err := data.loadData()
	if err != nil {
		log.Fatalf("failed to load data: %v\n", err)
	}

	// template
	err = rootTemplate.ExecuteTemplate(w, "base", struct {
		TableData string
	}{
		TableData: data.getDataHTML(),
	})

	if err != nil {
		http.Error(w, "failed to execute template", http.StatusNotFound)
		log.Println(err)
		return
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	r.ParseForm()

	err := data.add([]string{r.FormValue("name"), r.FormValue("description"), r.FormValue("count"), r.FormValue("date")})
	if err != nil {
		log.Printf("failed to add data: %v", err)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	r.ParseForm()
	n, err := strconv.Atoi(r.FormValue("number"))
	if err != nil {
		log.Println(err)
	}
	err = data.delete(n)
	if err != nil {
		log.Println(err)
	}
}

func undoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/undo" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	fmt.Println(data.cache)
	err := data.restore()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	// define flags
	customPath := flag.String("path", "", "specify a custom path to the data file")
	customPort := flag.Uint("port", 0, "specify a custom port (default: 8080)")
	flag.Parse()

	// use custom path if specified
	if *customPath != "" {
		data.contentPath = *customPath
	}

	// create data and cache file if they doesn't exist
	if _, err := os.Stat(data.contentPath); os.IsNotExist(err) {
		_, err := os.Create(data.contentPath)
		if err != nil {
			log.Fatalf("failed to create data file: %v\n", err)
		}
		// default table heading
		data.add([]string{"Name", "Description", "Count", "Date"})
	}
	if _, err := os.Stat(data.cachePath); os.IsNotExist(err) {
		_, err := os.Create(data.cachePath)
		if err != nil {
			log.Fatalf("failed to create cache file: %v\n", err)
		}
	}

	// server
	mux := http.NewServeMux()
	mux.Handle("/public/", public())
	mux.HandleFunc("/add", addHandler)
	mux.HandleFunc("/delete", deleteHandler)
	mux.HandleFunc("/undo", undoHandler)
	mux.HandleFunc("/", handleRoot)

	port := "8080"
	if *customPort != 0 {
		port = fmt.Sprint(*customPort)
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("inventory: starting server on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("inventory: couldn't start server: %v\n", err)
	}
}
