/*
 * File share api
 *
 * File share api.
 *
 * API version: 2.0.0
 * Contact: enponsba@gmail.com
 */

package main

import (
	"log"
	"net/http"
	"os"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myfilename/myrepo/go"
	//

	"github.com/gorilla/handlers"
	"github.com/ticquique/raspi-repo/backend/raspi"
)

func main() {

	raspi.NewDbConnection()

	os.Setenv("assetsDir", "assets")
	os.Setenv("ORIGIN_ALLOWED", "http://localhost:8080")

	if args := os.Args; len(args) == 2 && args[1] == "seed" {

		raspi.SeedDb()

	} else {

		log.Printf("Server started")

		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
		originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

		router := raspi.NewRouter()

		log.Fatal(http.ListenAndServe(":8081", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
	}
}
