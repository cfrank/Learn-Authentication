package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cfrank/auth.fun/api/authentication"
	"github.com/idawes/httptreemux"
)

func main() {
	const PORT string = ":8117"

	fmt.Printf("Starting web server on %s", PORT)

	// Define the base router
	router := httptreemux.New()

	// Routes under /auth
	authRoute := router.NewGroup("/auth")

	authRoute.POST("/signup", authentication.NewAuth)

	log.Fatal(http.ListenAndServe(PORT, router))
}
