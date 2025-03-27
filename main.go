package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Nivedithabp/receipt-processor/routes"
)

func main() {
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", router)
}