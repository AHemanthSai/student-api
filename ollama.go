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
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	s, exists := students[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Prompt construction
	prompt := "Summarize this student:\nName: " + s.Name + "\nEmail: " + s.Email + "\nAge: " + strconv.Itoa(s.Age)

	// Prepare the request payload
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	})

	// Use your ngrok-exposed Ollama server here
	ollamaURL := "https://5888-2405-201-c054-88da-5c4b-5d55-1c8c-c91e.ngrok-free.app/api/generate"


	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Error calling Ollama", http.StatusInternalServerError)
		log.Println("Error calling Ollama:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading Ollama response", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}

	summary, _ := result["response"].(string)
	json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}

