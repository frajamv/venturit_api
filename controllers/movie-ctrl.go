package controllers

import (
	"encoding/json"
	DML "moviesapi/db_operations"
	"moviesapi/models"
	"moviesapi/utils"
	"net/http"
	"strconv"
	"strings"
)

/**
* Fetch all the movies in the database with optional filters.
*
 */
func GetAllMovies(res http.ResponseWriter, req *http.Request) {
	title_filter := req.URL.Query()["title"]            // Title filters.
	id_filter := req.URL.Query()["id"]                  // Id filters.
	year_filter := req.URL.Query()["year"]              // Release year filters.
	rating_filter := req.URL.Query()["rating"]          // IMDB rating filters.
	genre_filter := req.URL.Query()["genre"]            // Genre filters.
	inclusive_filtering := req.URL.Query()["inclusive"] // Determine wether the filters will be inclusive or wether not.

	inclusive := "1"
	if inclusive_filtering != nil { // The filters are inclusive if this option is not provided.
		inclusive = inclusive_filtering[0]
	}

	// Get all movies with the selected filters (if there are).
	data, err := DML.FetchAllMovies(title_filter, id_filter, year_filter, rating_filter, genre_filter, inclusive)

	if err != nil {
		utils.SendErrorResponse(err, res)

	} else {
		if len(data) < 1 && title_filter != nil { // If no movies were found and there is a title filter, search in GOMDB API.
			title := ""
			if title_filter != nil {
				title = title_filter[0]
			}
			omdb_data, omdb_err := DML.FetchMovieByTitle_OMDB(title) // Get movie by title from GOMDB.

			if omdb_err == nil { // If the movie was found, insert into local database.
				genres := strings.Split(omdb_data.Genre, ", ") // Found movie genres.

				for _, genre_name := range genres {
					_, err := DML.FetchGenres(genre_name) // Fetch genre in local DB.
					if err != nil {
						DML.InsertGenre(genre_name) // If the genre is not on local DB, insert it.
					}
				}

				rating, _ := strconv.ParseFloat(omdb_data.ImdbRating, 64)                 // Found movie IMDB rating.
				movie_id, err := DML.InsertMovie(omdb_data.Title, omdb_data.Year, rating) // Add found GOMDB movie into local DB.

				if err != nil {
					utils.SendErrorResponse(omdb_err, res)
				}

				DML.InsertMovieGenres(movie_id, genres) // Relate inserted (or existing) genres to the newly created movie in local DB.

				utils.SendSuccessResponse(omdb_data, res) // Send the found movie in GOMDB.
			} else {
				utils.SendErrorResponse(omdb_err, res) // Send error message if movie was not found in GOMDB.
			}
		} else {
			utils.SendSuccessResponse(data, res) // If movie was found in local DB. Send it.
		}
	}
}

/**
* Update the selected movie rating and genres.
*
 */
func UpdateMovie(res http.ResponseWriter, req *http.Request) {
	var payload models.Movie
	_ = json.NewDecoder(req.Body).Decode(&payload)

	err := DML.UpdateMovieRating(payload.Id, payload.Rating)
	if err != nil {
		utils.SendErrorResponse(err, res)
	}

	for _, genre_name := range payload.Genres {
		_, err := DML.FetchGenres(genre_name) // Fetch genre in local DB.
		if err != nil {
			DML.InsertGenre(genre_name) // If the genre is not on local DB, insert it.
		}
	}

	err = DML.SetMovieGenres(payload.Id, payload.Genres)
	if err != nil {
		utils.SendErrorResponse(err, res)
	}
}
