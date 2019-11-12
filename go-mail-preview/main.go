package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattes/go-mail"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

var (
	listenFlag      string
	templateDirFlag string
)

func main() {
	flag.StringVar(&listenFlag, "listen", ":8000", "Listen on port")
	flag.StringVar(&templateDirFlag, "templates", "./templates", "Template directory")
	flag.Parse()

	files := mail.FilesFromLocalDir(templateDirFlag)

	mux := http.NewServeMux()

	mux.HandleFunc("/preview/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// parse templates
		t, err := mail.NewTemplates(files)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// parse name out of url
		name := strings.TrimPrefix(r.URL.Path, "/preview/")
		if name == "" {
			fmt.Fprint(w, "Error: Please visit /preview/xxx.html in your browser.")
			return
		}

		// render template with example data from name.example.yaml
		out, err := t.RenderWithExampleData(name)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// set content-type header
		switch filepath.Ext(name) {
		case "." + mail.HTMLExtension:
			w.Header().Set("content-type", "text/html")

		default:
			fallthrough
		case "." + mail.TextExtension:
			w.Header().Set("content-type", "text/plain")
		}

		log.Printf("Rendered %v in %v", name, time.Since(start))

		// write rendered template
		w.Write(out)
	})

	abs, err := filepath.Abs(templateDirFlag)
	if err != nil {
		abs = templateDirFlag
	}
	log.Printf("Listening at %v serving %v", listenFlag, abs)
	log.Fatal(http.ListenAndServe(listenFlag, mux))
}
