package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/students", CreateStudent).Methods("POST")
	r.HandleFunc("/students", GetStudents).Methods("GET")
	r.HandleFunc("/students/{id}", GetStudentByID).Methods("GET")
	r.HandleFunc("/students/{id}", UpdateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", DeleteStudent).Methods("DELETE")
	r.HandleFunc("/students/{id}/summary", GetStudentSummary).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local
	}

	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
