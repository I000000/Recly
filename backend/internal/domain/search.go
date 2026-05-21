package domain

type ItemDetail struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Type   string   `json:"type"` // "book" или "movie"
	Image  string   `json:"image"`
	Genres []string `json:"genres,omitempty"`

	Year        int     `json:"year,omitempty"`
	Rating      float64 `json:"rating,omitempty"`
	Description string  `json:"description,omitempty"`

	Authors  string `json:"authors,omitempty"`
	Director string `json:"director,omitempty"`
	Cast     string `json:"cast,omitempty"`
	Runtime  int    `json:"runtime,omitempty"`
}
