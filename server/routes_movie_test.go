package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/waldhalf/tplwebservergo/models"

	"github.com/stretchr/testify/assert"
)

type testStore struct {
	movieId int64
	movies [] *models.Movie
}

func (t testStore) Open() error {
	return nil
}

func (t testStore) Close() error {
	return nil
}
func (t testStore)GetMovies() ([]*models.Movie, error){
	return t.movies, nil
}

func (t testStore)GetMovieById(id int64) (*models.Movie, error){
	for _, m := range t.movies {
		if m.ID == id {
			return m, nil
 		}	
	}
	return nil, nil
}

func (t testStore)CreateMovie(m *models.Movie) (error){
	t.movieId++
	m.ID = t.movieId
	t.movies = append(t.movies, m)
	return nil
}

func (t testStore) FindUser(username string, password string)(bool, error){
	return true, nil
}
func TestMovieCreateUnit(t *testing.T){
	// Create Server with test DB
	srv := NewServer()
	srv.Store = &testStore{}

	// Prepare json body
	p := struct {
		Title string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration int `json:"duration"`
		TrailerUrl string `json:"trailer_url"`
	} {
		Title : "Inception",
		ReleaseDate: "2008-09-14",
		Duration: 148, 
		TrailerUrl: "http://url",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/api/movies", &buf)
	w := httptest.NewRecorder()
	srv.handleMovieCreate()(w,r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMovieCreateIntegration(t *testing.T){
	// Create Server with test DB
	srv := NewServer()
	srv.Store = &testStore{}

	// Prepare json body
	p := struct {
		Title string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration int `json:"duration"`
		TrailerUrl string `json:"trailer_url"`
	} {
		Title : "Inception",
		ReleaseDate: "2008-09-14",
		Duration: 148, 
		TrailerUrl: "http://url",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/api/movies", &buf)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzUxNDkzODUsImlhdCI6MTYzNTE0NTc4NSwidXNlcm5hbWUiOiJHb2xhbmcifQ.ptpGZkgBN62ikecBiNLfVzOaAm1s4Gkdvc4HP7lIrDM"
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w := httptest.NewRecorder()

	srv.serveHTTP(w,r)
	assert.Equal(t, http.StatusOK, w.Code)
}