package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/Nivedithabp/receipt-processor/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/Nivedithabp/receipt-processor/routes"
)

func main() {
	router := mux.NewRouter()

	// Register Routes
	routes.RegisterRoutes(router)

	// Serve Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	port := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
