package main

import (
	"log"
	"net/http"

	"github.com/aktsk/atgen/example/handlers"
)

// main function to boot up everything
func main() {
	log.Fatal(http.ListenAndServe(":8000", handlers.GetRouter()))
}
