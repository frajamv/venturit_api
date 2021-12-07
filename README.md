# venturit_api
 Venturit backend test - Movies REST API.


# Run REST API.
1. Insert .env file.
2. Open integrated terminal.
3. Execute "go get ." in integrated terminal to fetch and update all packages and dependencies.
4. Execute "go run main.go" in integrated terminal to execute de application.
5. Import Postman collection.
6. Run endpoints.


# Database
- The database is mounted online using AWS RDS service.
- A complete DDL SQL backup file for the local database is located in the file "db/database_ddl.sql".

# Methods

- GET movies/all?[filters]: Returns all movies from local database that match the selected filters. if "inclusive" parameter is provided with value of "1" or "true", the search will filter inclusively through all the movies in the local database. If filtering through title and the movie is not on local database, it will search through the GOMDB API for the movie with it's title, if found, it will be added to the local database for further searches.

- PATCH movies/update: Updates the movie which id corresponds to the id provided through the request body and updates it's IMDB rating and genres with the ones provided trough the request body.

# Developed by Francisco Javier Martínez Vargas as a backend test for Venturit on Dec, 2021.
