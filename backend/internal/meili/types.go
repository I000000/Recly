package meili

type rawHit struct {
	ID              interface{} `json:"id"`
	Title           string      `json:"title"`
	Type            string      `json:"type"`
	ImageURL        string      `json:"image_url"`
	Genres          interface{} `json:"genres"`
	AverageRating   interface{} `json:"average_rating"`
	VoteAverage     interface{} `json:"vote_average"`
	RatingsCount    interface{} `json:"ratings_count"`
	VoteCount       interface{} `json:"vote_count"`
	PublicationYear interface{} `json:"publication_year"`
	ReleaseDate     interface{} `json:"release_date"`
	Description     string      `json:"description"`
	Overview        string      `json:"overview"`
	Authors         string      `json:"authors"`
	Director        string      `json:"director"`
	Cast            string      `json:"cast"`
	Runtime         interface{} `json:"runtime"`
}
