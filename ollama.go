package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetStudentSummary(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	student, exists := students[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Build prompt
	prompt := fmt.Sprintf("Write a professional summary about a student named %s, aged %d, with email %s.", student.Name, student.Age, student.Email)

	// Create request body
	reqBody, _ := json.Marshal(map[string]string{
		"model":  "mistral", // or llama3 if it fits in your RAM
		"prompt": prompt,
	})

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "Error calling Ollama", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var summaryText string
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			break
		}
		if part, ok := chunk["response"].(string); ok {
			summaryText += part
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"summary": summaryText,
	})
}
