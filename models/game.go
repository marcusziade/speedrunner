package models

// Game represents a speedrun.com game in bulk mode
type Game struct {
	ID           string    `json:"id"`
	Names        GameNames `json:"names"`
	Abbreviation string    `json:"abbreviation"`
	Weblink      string    `json:"weblink"`
}

type GameNames struct {
	International string `json:"international"`
	Japanese      string `json:"japanese,omitempty"`
}
