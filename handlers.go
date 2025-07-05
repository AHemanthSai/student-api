package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var s Student
	json.NewDecoder(r.Body).Decode(&s)

	if s.ID == 0 || s.Name == "" || s.Email == "" || s.Age <= 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	students[s.ID] = s
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func GetStudents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(students)
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	s, exists := students[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, exists := students[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	var s Student
	json.NewDecoder(r.Body).Decode(&s)
	s.ID = id
	students[id] = s
	json.NewEncoder(w).Encode(s)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if _, exists := students[id]; !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	delete(students, id)
	w.WriteHeader(http.StatusNoContent)
}

