package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetStudentSummary(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	mu.RLock()
	s, exists := students[id]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// ✅ Mocked AI summary (no Ollama call)
	summary := "Here is a summary of the student:\n\n" +
		"* Name: " + s.Name + "\n" +
		"* Email: " + s.Email + "\n" +
		"* Age: " + strconv.Itoa(s.Age)

	json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}

