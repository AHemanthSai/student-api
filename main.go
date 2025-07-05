package main

import (
	"log"
	"net/http"

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

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}


