package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cfrank/auth.fun/api/auth"
	"github.com/cfrank/auth.fun/api/database"
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

	// Open connection to DB
	databaseError := database.Open()

	if databaseError != nil && database.MyDb.Alive {
		log.Fatal("Error: Problem establishing connection to database!")
	}

	// Defer closing the database
	defer database.Close()

	// Catch SIGTERM and close DB
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGTERM)
	go handleSigTerm(termChan)

	log.Fatal(http.ListenAndServe(PORT, router))
}

// Closes database
func handleSigTerm(c chan os.Signal) {
	<-c
	database.Close()
	os.Exit(1)
}
