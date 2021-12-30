package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type model struct {
	t *template.Template
	s *status
}

func (m *model) serveTemplate(w http.ResponseWriter, r *http.Request) {
	// Load templates and serve.
	m.t.Execute(w, m.s)
}

func (m *model) updater(addr string, delay time.Duration) {
	for {
		time.Sleep(delay)
		s, err := ping(addr)
		if err != nil {
			log.Printf("failed to ping minecraft server: %v\n", err)
		} else {
			m.s = s
		}
	}
}

func main() {
	// Read config file or fail.
	var c config
	cPath := "/etc/nilsu/mc.toml"
	if len(os.Args) > 1 {
		cPath = os.Args[1]
	}
	if err := c.Load(cPath); err != nil {
		log.Fatalf("failed to read config: %v\n", cPath)
	}
	log.Println("config loaded")

	// Create and setup model.
	var m model
	s, err := ping(c.MCAddress)
	if err != nil {
		log.Fatalf("failed to initially ping minecraft server: %v\n", err)
	}
	m.s = s

	// Update status every minute.
	go m.updater(c.MCAddress, time.Minute)

	t, err := template.ParseFiles(c.Template)
	if err != nil {
		log.Fatalf("failed to load and parse template: %v\n", err)
	}
	m.t = t

	// Serve or crash.
	http.HandleFunc("/", m.serveTemplate)
	log.Println("opening on:", c.Address)
	log.Fatal(http.ListenAndServe(c.Address, nil))
}
