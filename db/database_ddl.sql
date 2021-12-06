-- DDL documentation for movies REST API.

USE moviesdb; -- Database to connect.

DROP TABLE IF EXISTS movie, genre, movie_genre; -- Remove tables if they already exist.

CREATE TABLE movie ( -- Table for movies.
	movie_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR (255) NOT NULL,
    released_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    rating DOUBLE NOT NULL DEFAULT 0.0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE genre ( -- Table for movie genres.
	genre_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR (255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE movie_genre ( -- Table for relating movie(s) with genre(s).
	movie_genre_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    movie_id INT NOT NULL REFERENCES movie (movie_id),
    genre_id INT NOT NULL REFERENCES genre (genre_id),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


/**
* Francisco Javier Martínez Vargas
* Venturit backend test
* December/2021
*
**/