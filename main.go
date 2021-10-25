package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/waldhalf/tplwebservergo/server"
	"github.com/waldhalf/tplwebservergo/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error{
	srv := server.NewServer()
	
	// Database
	srv.Store = &store.DbStore{}
	err := srv.Store.Open()
	if err != nil {
		return err
	}
	http.HandleFunc("/", server.LogRequestMiddleware(srv.Router.ServeHTTP) )
	log.Printf("Serving HTTP on PORT 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}

	defer srv.Store.Close()
	return nil
}

