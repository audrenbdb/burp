package main

import (
	"burp"
	"burp/httpreq"
	"burp/repo"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", "localhost:8080", "server address")
	prod = flag.Bool("prod", false, "prod environment")
	psql = flag.Bool("psql", false, "persist with psql")
)

// lax CORS middleware. Do not use this in prod.
func devMW(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%-8s %s", r.Method, r.URL.Path))
		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")

			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			return
		}
		next.ServeHTTP(w, r)
	}
}

func main() {
	flag.Parse()

	r := getRepo()
	s := burp.NewService(r)

	handler := httpreq.NewHandler(s)
	if !*prod {
		log.Println("WARNING: dev mode enabled.")
		handler = devMW(handler)
	}

	log.Println("server starting: ", *addr)
	err := http.ListenAndServe(*addr, handler)
	if err != nil {
		log.Fatal(err)
	}

}

func getRepo() burp.Repo {
	if *psql {
		return repo.NewPSQL()
	}
	return repo.NewInMemory()
}
