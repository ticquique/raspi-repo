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
	//    sw "github.com/myname/myrepo/go"
	//

	"github.com/ticquique/raspi-repo/raspi"
)

func main() {

	raspi.NewDbConnection()

	if args := os.Args; len(args) == 2 {

		raspi.SeedDb(args[1])

	} else {

		log.Printf("Server started")

		router := raspi.NewRouter()

		log.Fatal(http.ListenAndServe(":8081", router))
	}
}
