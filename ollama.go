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

	// If running on Railway (no local Ollama), return mock summary
	if os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("OLLAMA_URL") == "" {
		mock := fmt.Sprintf("Mock summary for %s, age %d, email: %s.", s.Name, s.Age, s.Email)
		json.NewEncoder(w).Encode(map[string]string{"summary": mock})
		return
	}

	// If running locally, use Ollama
	prompt := "Summarize this student's profile:\n\n" +
		"Name: " + s.Name + "\n" +
		"Email: " + s.Email + "\n" +
		"Age: " + strconv.Itoa(s.Age)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  "mistral",
		"prompt": prompt,
		"stream": false,
	})
	if err != nil {
		http.Error(w, "Error encoding request body", http.StatusInternalServerError)
		return
	}

	ollamaURL := strings.TrimSpace(os.Getenv("OLLAMA_URL"))
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

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Error decoding Ollama response", http.StatusInternalServerError)
		return
	}

	if result.Response == "" {
		http.Error(w, "Empty Ollama response", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"summary": result.Response})
}

