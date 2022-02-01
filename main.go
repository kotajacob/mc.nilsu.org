package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type model struct {
	tmpl    *template.Template
	display display
}

type display struct {
	Status  *status
	Offline bool
	Mods    keySlice
	Carpets keySlice
}

func (m *model) serveTemplate(w http.ResponseWriter, r *http.Request) {
	// Load templates and serve.
	if err := m.tmpl.Execute(w, m.display); err != nil {
		log.Println(err)
	}
}

func (m *model) update(c config) {
	s, err := ping(c.MCAddress)
	if err != nil {
		log.Printf("failed to ping minecraft server: %v\n", err)
		m.display.Offline = true
	} else {
		m.display.Status = s
		m.display.Offline = false
	}

	m.display.Mods, err = parseKeyFile(c.ModList)
	if err != nil {
		log.Fatalf("failed parsing mod list: %v\n", err)
	}
	m.display.Carpets, err = parseKeyFile(c.CarpetList)
	if err != nil {
		log.Fatalf("failed parsing carpet list: %v\n", err)
	}
}

func (m *model) updater(c config, delay time.Duration) {
	for {
		time.Sleep(delay)
		m.update(c)
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
	m.update(c)

	// Update model every 5 minutes.
	go m.updater(c, 5*time.Minute)

	// Parse template and store in model.
	tmpl, err := template.ParseFiles(c.Template)
	if err != nil {
		log.Fatalf("failed to load and parse template: %v\n", err)
	}
	m.tmpl = tmpl

	// Serve or crash.
	http.HandleFunc("/", m.serveTemplate)
	log.Println("opening on:", c.Address)
	log.Fatal(http.ListenAndServe(c.Address, nil))
}
