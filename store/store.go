// L'interface store est le point d'entrée pour faire des requêtes SQL
// Cela nous permettra aussi d'écrire des TU plus facilement
package store

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/waldhalf/tplwebservergo/models"
)

type Storer interface {
	Open() error
	Close() error
	GetMovies() ([]*models.Movie, error)
	GetMovieById(id int64) (*models.Movie, error)
	CreateMovie(*models.Movie) error
	FindUser(username string, password string) (bool,error)
}

type DbStore struct {
	db *sqlx.DB
}

var schema = `
CREATE TABLE IF NOT EXISTS movies
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT,
	release_date TEXT,
	duration INTEGER,
	trailer_url TEXT
);

CREATE TABLE IF NOT EXISTS users
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	password TEXT
);
`

func (store *DbStore) Open() error{
	// On créé programmatiquement le folder database
	os.MkdirAll("database", 0777)
	// On créé / connect to DB
	db, err := sqlx.Connect("sqlite3", "database/db_go.db")
	if err != nil {
		return err
	}
	log.Println("Connected to db")
	store.db = db

	// On instancie la nouvelle table si elle n'existe pas
	db.MustExec(schema)

	return nil
}

func (store *DbStore) Close() error {
	return store.db.Close()
}

func (store *DbStore)GetMovies() ([]*models.Movie, error){
	var movies []*models.Movie
	err := store.db.Select(&movies, "SELECT * FROM movies")
	if err != nil {
		return movies, err
	}

	return movies, nil
}

func (store *DbStore)GetMovieById(id int64) (*models.Movie, error){
	var movie = &models.Movie{}
	err := store.db.Get(movie, "SELECT * FROM movies WHERE id=$1", id)
	if err != nil {
		return movie, err
	}

	return movie, nil
}

func (store *DbStore)CreateMovie(m *models.Movie) (error){
	res , err := store.db.Exec("INSERT INTO movies (title, release_date, duration, trailer_url) VALUES(?,?,?,?)", 
	m.Title,
	m.ReleaseDate, 
	m.Duration,
	m.TrailerUrl)
	if err != nil {
		return err
	}
	m.ID, err = res.LastInsertId()
	return err
}

func (store *DbStore)FindUser(username, password string) (bool, error){
	var count int
	err :=store.db.Get(&count, "SELECT count(id) FROM users where username=$1 AND password=$2", username, password)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}
