package db_operations

import (
	DB "moviesapi/db"
	"moviesapi/models"
	"os"
	"strconv"
	"strings"

	"github.com/eefret/gomdb"
)

/**
* Creates a SQL statement to get all the movies (or with a custom filtering) and returns an array of movies.
 */
func FetchAllMovies(title_filter []string, id_filter []string, year_filter []string, rating_filter []string, genre_filter []string, inclusive_filtering string) ([]models.Movie, error) {
	movie_result := []models.Movie{}
	bd, err := DB.GetDBconnection()
	if err != nil {
		return movie_result, err
	}

	// Retrieves the movies with the related genres.
	query := "SELECT "
	query += "Movie.movie_id AS id, Movie.title, YEAR(released_date) AS released_year, Movie.rating, "
	query += "GROUP_CONCAT(Genre.name SEPARATOR ', ') AS genres "
	query += "FROM movie_genre AS MG "
	query += "JOIN movie AS Movie USING (movie_id) "
	query += "JOIN genre AS Genre USING (genre_id) "

	query = buildQueryFilters(query, title_filter, id_filter, year_filter, rating_filter, genre_filter, inclusive_filtering)

	query += "GROUP BY Movie.movie_id, Movie.title, Movie.rating;"

	rows, err := bd.Query(query)
	if err != nil {
		return movie_result, err
	}

	for rows.Next() { // For each movie row.
		var movie models.Movie
		var genres string

		err = rows.Scan(&movie.Id, &movie.Title, &movie.Released_year, &movie.Rating, &genres)
		movie.Genres = strings.Split(genres, ", ") // Convert genres string into array.
		if err != nil {
			return movie_result, err
		}
		movie_result = append(movie_result, movie)
	}

	return movie_result, nil
}

/**
* Sets up SQL statements for every filtering array and determines wether these filters will be inclusive or exclusive.
 */
func buildQueryFilters(query string, title_filter []string, id_filter []string, year_filter []string, rating_filter []string, genre_filter []string, inclusive_filtering string) string {
	by_title := len(title_filter) > 0
	by_id := len(id_filter) > 0
	by_year := len(year_filter) > 0
	by_rating := len(rating_filter) > 0
	by_genre := len(genre_filter) > 0

	if by_title || by_id || by_year || by_rating || by_genre { // If filtering movies. Must: Define inclusivity, filter operator and filter values.

		inclusive, _ := strconv.ParseBool(inclusive_filtering)
		query += "WHERE " + strconv.FormatBool(!inclusive) + " " // Depending on operator so it won't disturb the results.

		operator := "AND" // Exclusive filtering.
		if inclusive {    // Determine wether the filters will be inclusive or exclusive.
			operator = "OR" // Inclusive filtering
		}

		if by_title { // Movies which title is within the title filters.
			query += operator + " Movie.title IN ('" + strings.Join(title_filter, "', '") + "') "
		}

		if by_id { // Movies which id is within the id filters.
			query += operator + " Movie.movie_id IN (" + strings.Join(id_filter, ", ") + ") "
		}

		if by_year { // Movies which year is within the year filters.
			query += operator + " YEAR(Movie.released_date) IN (" + strings.Join(year_filter, ", ") + ") "
		}

		if by_rating { // Movies which rating is higher or lower than the rating filters.
			query += operator + " Movie.rating <> (" + strings.Join(rating_filter, ", ") + ") "
		}

		if by_genre { // Movies which genre is within the genre filters.
			query += operator + " Movie.title IN ('" + strings.Join(genre_filter, "', '") + "') "
		}
	}

	return query
}

/**
* Retrieves movies from the OMDB API within the selected title, if exist.
 */
func FetchMovieByTitle_OMDB(title string) (*gomdb.MovieResult, error) {
	OMDB_API_KEY := os.Getenv("OMDB_API_KEY")

	api := gomdb.Init(OMDB_API_KEY)
	query := &gomdb.QueryData{Title: title, SearchType: gomdb.MovieSearch}
	res, err := api.MovieByTitle(query)

	if err != nil {
		return nil, err
	}
	return res, nil
}

/**
* Creates a new movie with the specified parameters.
 */
func InsertMovie(title string, release_year string, rating float64) (int64, error) {
	bd, err := DB.GetDBconnection()
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO movie (title, released_date, rating) VALUES (?, ?, ?)"
	result, err := bd.Exec(query, title, release_year+"-01-01", rating)
	if err != nil {
		return 0, err
	}

	movie_id, err := result.LastInsertId() // Get inserted object id for later genres assignment.
	if err != nil {
		return 0, err
	}

	return movie_id, nil
}

/**
* Retrieves all the genres in the local DB.
**/
func FetchGenres(genre_name string) (models.Genre, error) {
	var genre models.Genre
	bd, err := DB.GetDBconnection()
	if err != nil {
		return genre, err
	}

	query := "SELECT genre_id AS id, name FROM genre WHERE name = '" + genre_name + "' LIMIT 1"
	result := bd.QueryRow(query)
	err = result.Scan(&genre.Id, &genre.Name)

	if err != nil {
		return genre, err
	}

	return genre, nil
}

/**
* Insert a new genre in the local DB.
**/
func InsertGenre(genre_name string) error {
	bd, err := DB.GetDBconnection()
	if err != nil {
		return err
	}

	query := "INSERT INTO genre (name) VALUES ('" + genre_name + "')"
	_, err = bd.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

/**
* Insert a new row in the movie_genre table to assign X number of genres to a movie.
**/
func InsertMovieGenres(movie_id int64, genre_names []string) error {
	bd, err := DB.GetDBconnection()
	if err != nil {
		return err
	}

	query := "INSERT INTO movie_genre (movie_id, genre_id) (SELECT " + strconv.Itoa(int(movie_id)) + ", genre_id FROM genre WHERE name IN ('" + strings.Join(genre_names, "', '") + "'))"
	_, err = bd.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

/**
* Updates a movie IMDB rating in the local DB.
**/
func UpdateMovieRating(movie_id int64, new_rating float64) error {
	bd, err := DB.GetDBconnection()
	if err != nil {
		return err
	}

	query := "UPDATE movie SET rating = ? WHERE movie_id = ?"
	_, err = bd.Exec(query, new_rating, movie_id)
	if err != nil {
		return err
	}

	return nil
}

/**
* Deletes all genres from the selected movie and assigns all the new genres set as parameters.
**/
func SetMovieGenres(movie_id int64, genres []string) error {
	bd, err := DB.GetDBconnection()
	if err != nil {
		return err
	}

	query := "DELETE FROM movie_genre WHERE movie_id = ?" // Deletes all the movie genres.
	_, err = bd.Exec(query, movie_id)
	if err != nil {
		return err
	}

	// Relates all the new genres to the selected movie.
	query = "INSERT INTO movie_genre (movie_id, genre_id) (SELECT " + strconv.Itoa(int(movie_id)) + ", genre_id FROM genre WHERE name IN ('" + strings.Join(genres, "', '") + "'))"
	_, err = bd.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
