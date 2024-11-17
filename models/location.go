package models

type Location struct {
	Country Country `json:"country"`
	Region  *Region `json:"region,omitempty"`
}

type Country struct {
	Code  string       `json:"code"`
	Names CountryNames `json:"names"`
}

type CountryNames struct {
	International string `json:"international"`
	Japanese      string `json:"japanese,omitempty"`
}

type Region struct {
	Code  string      `json:"code"`
	Names RegionNames `json:"names"`
}

type RegionNames struct {
	International string `json:"international"`
	Japanese      string `json:"japanese,omitempty"`
}
