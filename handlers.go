package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if s.ID == 0 || s.Name == "" || s.Email == "" || s.Age <= 0 {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	students[s.ID] = s
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func GetStudents(w http.ResponseWriter, r *http.Request) {
	var all []Student
	for _, s := range students {
		all = append(all, s)
	}
	json.NewEncoder(w).Encode(all)
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
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
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

