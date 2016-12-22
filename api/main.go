package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cfrank/auth.fun/api/auth"
	"github.com/idawes/httptreemux"
)

func main() {
	const PORT string = ":8117"

	fmt.Printf("Starting web server on %s\n", PORT)

	// Define the base router
	router := httptreemux.New()

	// Routes under /auth
	authRoute := router.NewGroup("/auth")

	authRoute.POST("/signup", auth.NewAuth)

	log.Fatal(http.ListenAndServe(PORT, router))
}
