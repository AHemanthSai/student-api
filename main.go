package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Student API is running! 🚀"))
	})

	r.HandleFunc("/students", CreateStudent).Methods("POST")
	r.HandleFunc("/students", GetStudents).Methods("GET")
	r.HandleFunc("/students/{id}", GetStudentByID).Methods("GET")
	r.HandleFunc("/students/{id}", UpdateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", DeleteStudent).Methods("DELETE")
	r.HandleFunc("/students/{id}/summary", GetStudentSummary).Methods("GET")

	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
