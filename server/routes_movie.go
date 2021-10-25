package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/waldhalf/tplwebservergo/models"
)

type jsonMovie struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Duration int `json:"duration"`
	TrailerUrl string `json:"trailer_url"`
}

type jsonError struct {
	Message string `json:"message"`
}


func(s *server)handleMovieList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		movies, err := s.Store.GetMovies()
		if err != nil{
			log.Printf("Cannot load movies. err : %v\n", err)
			// TODO handle response to the client
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}
		// on convertit nos films en format JSON
		var resp = make([]jsonMovie, len(movies))
		for i, m := range movies {
			resp[i] = mapMovieToJson(m)
		}
		// TODO response JSON format
		s.respond(w, r, resp, http.StatusOK)
	}
}

func(s *server)handleMovieDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Cannot parse id to int. err : %v", err)
			s.respond(
				w,
				r,
				mapErrorToJson("Cannot parse id to int"),
				http.StatusBadRequest)
			return
		}

		m, err := s.Store.GetMovieById(id)
		if err != nil {
			log.Printf("Cannot getmovie by id. err : %v", err)
			s.respond(
				w, 
				r, 
				mapErrorToJson("Cannot getmovie by id"),
				http.StatusInternalServerError)
			return
		}

		var resp = mapMovieToJson(m)
		s.respond(w, r, resp, http.StatusOK)
	}
}

func(s *server)handleMovieCreate() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration int `json:"duration"`
		TrailerUrl string `json:"trailer_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Can't parse movie body. err : %v", err)
			s.respond(
			w, 
			r, 
			mapErrorToJson("Can't parse movie body"),
			http.StatusBadRequest)
			return 
		}
		// Create a movie
		m := &models.Movie{
			ID : 0,
			Title : req.Title,
			ReleaseDate: req.ReleaseDate,
			Duration : req.Duration,
			TrailerUrl : req.TrailerUrl,
		}

		// Store movie in DB
		err = s.Store.CreateMovie(m)
		if err != nil {
			log.Printf("Can't create movie in DB. err : %v", err)
			s.respond(
			w, 
			r, 
			mapErrorToJson("Can't create movie in DB"),
			http.StatusInternalServerError)
			return 
		}
		var resp = mapMovieToJson(m)
		s.respond(w, r, resp, http.StatusOK)
	}
}


func mapMovieToJson(m *models.Movie) jsonMovie{
	return jsonMovie{
		ID : m.ID,
		Title : m.Title,
		ReleaseDate : m.ReleaseDate,
		Duration: m.Duration,
		TrailerUrl:m.TrailerUrl,
	}
}

func mapErrorToJson(err string) jsonError{
	return jsonError {
		Message: err,
	}
}
