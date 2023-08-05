package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("../model"))

	log.Fatal(http.ListenAndServe(":9000", fs))
}
