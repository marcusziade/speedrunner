package ui

import (
	"fmt"

	"speedrunner/api"
	"speedrunner/models"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type SearchHandler struct {
	client      *api.Client
	searchEntry *gtk.SearchEntry
	searchType  *gtk.ComboBoxText
	spinner     *gtk.Spinner
	resultsList *gtk.ListBox
	statusFunc  func(string)
}

func NewSearchHandler(client *api.Client, updateStatus func(string)) (*SearchHandler, error) {
	// Create search type combo
	searchType, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	searchType.Append("games", "Games")
	searchType.Append("users", "Users")
	searchType.SetActive(0)

	// Create search entry
	searchEntry, err := gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	searchEntry.SetPlaceholderText("Search games or users...")

	// Create spinner
	spinner, err := gtk.SpinnerNew()
	if err != nil {
		return nil, err
	}

	// Create results list
	resultsList, err := gtk.ListBoxNew()
	if err != nil {
		return nil, err
	}
	resultsList.SetSelectionMode(gtk.SELECTION_SINGLE)

	handler := &SearchHandler{
		client:      client,
		searchEntry: searchEntry,
		searchType:  searchType,
		spinner:     spinner,
		resultsList: resultsList,
		statusFunc:  updateStatus,
	}

	// Connect signals
	searchEntry.Connect("activate", handler.onSearch)
	searchType.Connect("changed", handler.onSearchTypeChanged)
	resultsList.Connect("row-activated", handler.onResultActivated)

	return handler, nil
}

func (s *SearchHandler) GetSearchBox() (*gtk.Box, error) {
	searchBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, err
	}
	searchBox.SetMarginStart(10)
	searchBox.SetMarginEnd(10)
	searchBox.SetMarginTop(10)

	searchBox.PackStart(s.searchType, false, false, 0)
	searchBox.PackStart(s.searchEntry, true, true, 0)
	searchBox.PackStart(s.spinner, false, false, 0)

	return searchBox, nil
}

func (s *SearchHandler) GetResultsList() *gtk.ListBox {
	return s.resultsList
}

func (s *SearchHandler) onSearch() {
	searchText, err := s.searchEntry.GetText()
	if err != nil {
		return
	}
	if searchText == "" {
		return
	}

	// Start spinner and disable search
	s.spinner.Start()
	s.searchEntry.SetSensitive(false)

	// Clear previous results
	s.clearResults()

	// Get search type
	searchType := s.searchType.GetActiveID()

	// Perform search in a goroutine
	go func() {
		var results interface{}
		var err error

		switch searchType {
		case "games":
			results, err = s.client.SearchGames(searchText)
		case "users":
			results, err = s.client.SearchUsers(searchText)
		}

		// Update UI in main thread
		glib.IdleAdd(func() {
			s.spinner.Stop()
			s.searchEntry.SetSensitive(true)

			if err != nil {
				// TODO: Add error handling
				return
			}

			s.displayResults(results, searchType)
		})
	}()
}

func (s *SearchHandler) clearResults() {
	for {
		row := s.resultsList.GetRowAtIndex(0)
		if row == nil {
			break
		}
		s.resultsList.Remove(row)
	}
}

func (s *SearchHandler) displayResults(results interface{}, resultType string) {
	count := 0

	switch resultType {
	case "games":
		if games, ok := results.([]models.Game); ok {
			for _, game := range games {
				row := s.createGameRow(game)
				s.resultsList.Add(row)
			}
			count = len(games)
		}
	case "users":
		if users, ok := results.([]models.User); ok {
			for _, user := range users {
				row := s.createUserRow(user)
				s.resultsList.Add(row)
			}
			count = len(users)
		}
	}

	s.resultsList.ShowAll()
	if s.statusFunc != nil {
		s.statusFunc(fmt.Sprintf("Found %d results", count))
	}
}

func (s *SearchHandler) createGameRow(game models.Game) *gtk.ListBoxRow {
	row, _ := gtk.ListBoxRowNew()
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	box.SetMarginStart(10)
	box.SetMarginEnd(10)
	box.SetMarginTop(5)
	box.SetMarginBottom(5)

	label, _ := gtk.LabelNew("")
	label.SetMarkup(fmt.Sprintf("<b>%s</b>\n<small>%s</small>",
		game.Names.International, game.Abbreviation))
	label.SetHAlign(gtk.ALIGN_START)
	box.PackStart(label, true, true, 0)

	linkButton, _ := gtk.LinkButtonNew(game.Weblink)
	box.PackEnd(linkButton, false, false, 0)

	row.Add(box)
	return row
}

func (s *SearchHandler) createUserRow(user models.User) *gtk.ListBoxRow {
	row, _ := gtk.ListBoxRowNew()
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	box.SetMarginStart(10)
	box.SetMarginEnd(10)
	box.SetMarginTop(5)
	box.SetMarginBottom(5)

	label, _ := gtk.LabelNew("")
	locationText := ""
	if user.Location != nil && user.Location.Country.Names.International != "" {
		locationText = fmt.Sprintf("\n<small>%s</small>",
			user.Location.Country.Names.International)
	}
	label.SetMarkup(fmt.Sprintf("<b>%s</b>%s",
		user.Names.International, locationText))
	label.SetHAlign(gtk.ALIGN_START)
	box.PackStart(label, true, true, 0)

	linkButton, _ := gtk.LinkButtonNew(user.Weblink)
	box.PackEnd(linkButton, false, false, 0)

	row.Add(box)
	return row
}

func (s *SearchHandler) onSearchTypeChanged() {
	searchType := s.searchType.GetActiveID()
	switch searchType {
	case "games":
		s.searchEntry.SetPlaceholderText("Search for games...")
	case "users":
		s.searchEntry.SetPlaceholderText("Search for users...")
	}
	s.clearResults()
}

func (s *SearchHandler) onResultActivated(list *gtk.ListBox, row *gtk.ListBoxRow) {
	// Handle row selection/activation
	// This can be expanded to show detailed views
}
