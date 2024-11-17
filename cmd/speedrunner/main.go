package main

import (
	"log"
	"os"

	"speedrunner/api"
	"speedrunner/ui"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	log.Println("Starting application...")

	// Initialize GTK
	gtk.Init(&os.Args)
	log.Println("GTK initialized")

	// Create API client
	client := api.NewClient()
	log.Println("API client created")

	// Create main window
	log.Println("Creating main window...")
	win, err := ui.NewMainWindow(client)
	if err != nil {
		log.Fatal("Error creating window:", err)
	}
	log.Println("Main window created")

	// Show window and start main loop
	win.ShowAll()
	log.Println("Window shown, starting main loop")
	gtk.Main()
}
