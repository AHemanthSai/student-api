package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetStudentSummary(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	s, exists := students[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Create prompt
	prompt := "Summarize this student's profile:\n\n" +
		"Name: " + s.Name + "\n" +
		"Email: " + s.Email + "\n" +
		"Age: " + strconv.Itoa(s.Age)

	// Ollama expects this structure
	requestBody, err := json.Marshal(map[string]string{
		"model":  "llama3",
		"prompt": prompt,
	})
	if err != nil {
		http.Error(w, "Failed to encode prompt", http.StatusInternalServerError)
		return
	}

	// Replace with your ngrok tunnel URL (port 11434 exposed via ngrok)
	ollamaURL := "https://0654-2405-201-c054-88da-5c4b-5d55-1c8c-c91e.ngrok-free.app/api/generate"

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error calling Ollama:", err)
		http.Error(w, "Error calling Ollama", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Ollama response", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Invalid JSON from Ollama", http.StatusInternalServerError)
		return
	}

	summary, ok := result["response"].(string)
	if !ok {
		http.Error(w, "Ollama response missing", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}


