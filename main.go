package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/tusupov/gousers/handle"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
			log.Fatalf("main() panic: %v \ndebug stack: %s", err, debug.Stack())
		}
	}()

	var addr = flag.String("h", ":8080", "Host address")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/user", handle.UserNew).Methods("POST")
	r.HandleFunc("/user/{id:[0-9]+}", handle.User).Methods("GET")
	r.HandleFunc("/user/transfer", handle.UserAmountTransfer).Methods("POST")

	log.Printf("Listening [%s] ...\n", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatalln(err)
	}

}
