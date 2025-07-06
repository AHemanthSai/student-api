package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GetStudentSummary(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	s, exists := students[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	prompt := "Summarize this student's profile:\n\n" +
		"Name: " + s.Name + "\n" +
		"Email: " + s.Email + "\n" +
		"Age: " + strconv.Itoa(s.Age)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	})

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		http.Error(w, "OLLAMA_URL environment variable not set", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Error calling Ollama: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading Ollama response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Try strict decode first
	var result struct {
		Response string `json:"response"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil || result.Response == "" {
		// Fallback to extracting plain text if response is not JSON
		text := strings.TrimSpace(string(body))
		if text == "" {
			http.Error(w, "Empty Ollama response", http.StatusInternalServerError)
			return
		}
		// Send raw string as summary
		json.NewEncoder(w).Encode(map[string]string{"summary": text})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"summary": result.Response})
}


