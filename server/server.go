package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/waldhalf/tplwebservergo/store"
)

const JWT_APP_KEY = "testtest"

type server struct {
	Router *mux.Router
	Store store.Storer
}

func NewServer() *server{
	s := &server{
		Router : mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *server) serveHTTP(w http.ResponseWriter, r *http.Request){
	s.Router.ServeHTTP(w, r)
}

func (s *server) respond(w http.ResponseWriter, _ *http.Request, data interface{}, status int){
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	err :=json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Cannot format json. err : %v", err)
	}
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}