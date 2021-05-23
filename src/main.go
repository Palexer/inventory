package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type textConfirmation struct {
	Text string `json:"Text"`
}

// rootTemplate references the specified rootTemplate and caches the parsed results
// to help speed up response times.
var (
	rootTemplate = template.Must(template.ParseFiles("./templates/base.html", "./templates/index.html"))
	limiter      = NewIPRateLimiter(1, 5)
	version      = "0.8"
)

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

	// remove cache file when reloading page
	if err := os.Remove(data.cachePath); err != nil {
		log.Printf("failed to remove cache file: %v\n", err)
	}

	// load data
	err := data.loadData()
	if err != nil {
		log.Fatalf("failed to load data: %v\n", err)
	}

	// template
	err = rootTemplate.ExecuteTemplate(w, "base", struct {
		TableData string
		Heading   string
	}{
		TableData: data.getDataHTML(),
		Heading:   data.heading,
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

	t, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		log.Printf("failed to parse time: %v\n", err)
	}

	err = data.add([]string{r.FormValue("name"), r.FormValue("description"), r.FormValue("count"), t.Format("02.01.2006")})
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

	decoder := json.NewDecoder(r.Body)
	var text textConfirmation
	err := decoder.Decode(&text)
	if err != nil {
		log.Printf("failed to decode json data: %v\n", err)
	}

	n, err := strconv.Atoi(text.Text)
	if err != nil {
		log.Println(err)
		return
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
	err := data.restore()
	if err != nil {
		log.Println(err)
	}
}

func deleteAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/deleteall" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var text textConfirmation
	err := decoder.Decode(&text)
	if err != nil {
		log.Printf("failed to decode json data: %v\n", err)
	}

	if text.Text == data.key {
		// delete file (whole table) and create backup
		err = data.deleteAllAndBackUp()
		if err != nil {
			log.Printf("failed to delete file and back up: %v\n", err)
		}
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/edit" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	r.ParseForm()
	index, err := strconv.Atoi(r.FormValue("index"))
	if err != nil {
		log.Printf("failed to convert index to number: %v\n", err)
		return
	}

	t, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		log.Printf("failed to parse time: %v\n", err)
	}

	err = data.edit(index, []string{r.FormValue("name"), r.FormValue("description"), r.FormValue("count"), t.Format("02.01.2006")})
	if err != nil {
		log.Printf("failed to add data: %v", err)
	}
}

func main() {
	// define flags
	key := flag.String("key", "", "specify a custom key to delete the table")
	heading := flag.String("heading", "", "specify a custom table heading")
	customPort := flag.Uint("port", 0, "specify a custom port (default: 8080)")
	flag.Parse()

	for _, v := range flag.Args() {
		if strings.ToLower(v) == "version" {
			fmt.Printf("inventory version %s\n", version)
			os.Exit(0)
		}
	}

	if *key != "" {
		data.key = *key
	} else {
		data.key = "Inventory"
	}

	if *heading != "" {
		data.heading = *heading
	} else {
		data.heading = "Inventory"
	}

	port := "8080"
	if *customPort != 0 {
		port = fmt.Sprint(*customPort)
	}

	// server
	mux := http.NewServeMux()
	mux.Handle("/public/", public())
	mux.HandleFunc("/add", addHandler)
	mux.HandleFunc("/delete", deleteHandler)
	mux.HandleFunc("/undo", undoHandler)
	mux.HandleFunc("/deleteall", deleteAllHandler)
	mux.HandleFunc("/edit", editHandler)
	mux.HandleFunc("/", handleRoot)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      limitMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("inventory: starting server on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("inventory: couldn't start server: %v\n", err)
	}
}
