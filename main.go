package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
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

// pollUpdate pings a minecraft server to get a player count and version.
func (m *model) pollUpdate(c config) {
	s, err := ping(c.MCAddress)
	if err != nil {
		log.Printf("failed to ping minecraft server: %v\n", err)
		m.display.Offline = true
	} else {
		m.display.Status = s
		m.display.Offline = false
	}
}

// pollUpdater runs pollUpdate repeatedly with a delay.
func (m *model) pollUpdater(c config, delay time.Duration) {
	for {
		time.Sleep(delay)
		m.pollUpdate(c)
	}
}

// watchUpdate gets the mod and carpet rule lists.
func (m *model) watchUpdate(c config) {
	var err error
	m.display.Mods, err = parseKeyFile(c.ModList)
	if err != nil {
		log.Printf("failed parsing mod list: %v\n", err)
	}
	m.display.Carpets, err = parseKeyFile(c.CarpetList)
	if err != nil {
		log.Printf("failed parsing carpet list: %v\n", err)
	}
}

// watchUpdater reacts to fsnotify events by calling watchUpdate.
func (m *model) watchUpdater(c config, watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
				m.watchUpdate(c)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
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
	m.pollUpdate(c)
	m.watchUpdate(c)

	// Setup file watcher for mod and carpet rule lists.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(c.ModList)
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add(c.CarpetList)
	if err != nil {
		log.Fatal(err)
	}

	go m.watchUpdater(c, watcher)

	// Poll the MC server every 5 minutes.
	go m.pollUpdater(c, 5*time.Minute)

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
