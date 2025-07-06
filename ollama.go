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
		"model":  "mistral", // Or "llama3" if available
		"prompt": prompt,
	})

	// ❗ UPDATE THIS: Use your current ngrok public URL here
	ollamaURL := "https://07be-2405-201-c054-88da-5c4b-5d55-1c8c-c91e.ngrok-free.app/api/generate"

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(reqBody))
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"summary": summaryText,
	})
}

