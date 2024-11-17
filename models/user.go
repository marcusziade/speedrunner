package models

// User represents a speedrun.com user
type User struct {
	ID       string    `json:"id"`
	Names    UserNames `json:"names"`
	Weblink  string    `json:"weblink"`
	Location *Location `json:"location,omitempty"`
}

type UserNames struct {
	International string `json:"international"`
	Japanese      string `json:"japanese,omitempty"`
}
