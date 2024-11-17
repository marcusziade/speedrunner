package ui

import (
	"speedrunner/api"

	"github.com/gotk3/gotk3/gtk"
)

type MainWindow struct {
	*gtk.Window
	client        *api.Client
	searchHandler *SearchHandler
	statusBar     *gtk.Statusbar
}

func NewMainWindow(client *api.Client) (*gtk.Window, error) {
	// Create window
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}

	// Create main window struct
	mainWin := &MainWindow{
		Window: win,
		client: client,
	}

	// Set up main window
	win.SetTitle("Speedrun Browser")
	win.SetDefaultSize(800, 600)
	win.Connect("destroy", gtk.MainQuit)

	// Create main layout
	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}
	win.Add(mainBox)

	// Create header bar
	header, err := createHeaderBar()
	if err != nil {
		return nil, err
	}
	win.SetTitlebar(header)

	// Create status bar first so we can pass its update function
	statusBar, err := gtk.StatusbarNew()
	if err != nil {
		return nil, err
	}
	mainWin.statusBar = statusBar

	// Create search handler with status update function
	searchHandler, err := NewSearchHandler(client, func(msg string) {
		mainWin.UpdateStatus(msg)
	})
	if err != nil {
		return nil, err
	}
	mainWin.searchHandler = searchHandler

	// Add search box
	searchBox, err := searchHandler.GetSearchBox()
	if err != nil {
		return nil, err
	}
	mainBox.PackStart(searchBox, false, false, 0)

	// Create scrolled window for results
	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	scroll.SetVExpand(true)
	mainBox.PackStart(scroll, true, true, 0)

	// Add results list to scroll window
	scroll.Add(searchHandler.GetResultsList())

	// Add status bar
	mainBox.PackStart(statusBar, false, false, 0)

	return win, nil
}

func createHeaderBar() (*gtk.HeaderBar, error) {
	header, err := gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}

	header.SetShowCloseButton(true)
	header.SetTitle("Speedrun Browser")
	header.SetSubtitle("Search games and runners")

	return header, nil
}

func (w *MainWindow) UpdateStatus(message string) {
	w.statusBar.Push(w.statusBar.GetContextId("results"), message)
}
