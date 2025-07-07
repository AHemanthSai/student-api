package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var once sync.Once

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if s.ID == 0 || s.Name == "" || s.Email == "" || s.Age <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := students[s.ID]; exists {
		http.Error(w, "Student with this ID already exists", http.StatusBadRequest)
		return
	}

	students[s.ID] = s
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func GetStudents(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	var all []Student
	for _, s := range students {
		all = append(all, s)
	}
	json.NewEncoder(w).Encode(all)
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	s, exists := students[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

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

	if s.Name == "" || s.Email == "" || s.Age <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	s.ID = id
	students[id] = s
	json.NewEncoder(w).Encode(s)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := students[id]; !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	delete(students, id)
	w.WriteHeader(http.StatusNoContent)
}


