package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	httpAddr := flag.String("http", ":9888", "Server bind address")
	dir := flag.String("dir", "", "Directory path to serve")
	flag.Parse()

	err := http.ListenAndServe(*httpAddr, http.FileServer(http.Dir(*dir)))
	if err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
