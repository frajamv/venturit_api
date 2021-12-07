package models

// Movie object from DB.
type Movie struct {
	Id            int64    `json:"id"`
	Title         string   `json:"title"`
	Released_year int      `json:"released_year"`
	Rating        float64  `json:"rating"`
	Genres        []string `json:"genres"`
}
