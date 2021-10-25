package server


func (s *server) routes() {
	s.Router.HandleFunc("/", s.HandleIndex()).Methods("GET")
	s.Router.HandleFunc("/api/token", s.HandleTokenCreate()).Methods("POST")
	s.Router.HandleFunc("/api/movies/{id:[0-9]+}", s.handleMovieDetail()).Methods("GET")
	s.Router.HandleFunc("/api/movies", s.LoggedOnly(s.handleMovieList())).Methods("GET")
	s.Router.HandleFunc("/api/movies", s.LoggedOnly(s.handleMovieCreate())).Methods("POST")
}