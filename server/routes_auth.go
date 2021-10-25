package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

func (s *server) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Welcome to template webserver GO")
	}
}

func (s *server) HandleTokenCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`

	}

	type response struct {
		Token string `json:"token"`
	}

	type responseEror struct {
		Error string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request){
		// Parsing login body
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			msg := fmt.Sprintf("Can't parse login body. err : %v", err)
			log.Println(msg)
			s.respond(w, r, responseEror{
				Error: msg, 
			}, http.StatusBadRequest)
			return
		}

		// Check Credentials
		found, err :=s.Store.FindUser(req.Username, req.Password)
		if err != nil {
			s.respond(w,r, responseEror{
				Error : "Can't find user",
			}, http.StatusInternalServerError)
			return
		}
		if !found {
			s.respond(w,r, responseEror{
				Error : "Invalid Credentials",
			}, http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			"iat": time.Now().Unix(),
		})
		// Generate Token
		tokenStr, err := token.SignedString([]byte(JWT_APP_KEY))
		if err != nil {
			msg := fmt.Sprintf("Can't generate jwt. err : %v", err)
			s.respond(w,r, responseEror{
				Error : msg,
			}, http.StatusInternalServerError)
		}

		s.respond(w,r, response{
				Token : tokenStr,
			}, http.StatusOK)
	}
}