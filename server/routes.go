package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	router.HandleFunc("/login", LoginHandler).Methods("GET")
	router.HandleFunc("/logout", LogoutHandler).Methods("GET")

	router.HandleFunc("/upload", FileUploadHandler).Methods("POST")

	router.HandleFunc("/status/{userID}", UploadStatusHandler).Methods("GET")

	httpAddr := os.Getenv("HTTP_ADDR")
	log.Printf("Server starting at %s", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, router))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
}

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
}

func UploadStatusHandler(w http.ResponseWriter, r *http.Request) {
}