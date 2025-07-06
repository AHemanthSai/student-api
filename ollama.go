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

	prompt := "Summarize this student's profile:\n\n" +
		"Name: " + s.Name + "\n" +
		"Email: " + s.Email + "\n" +
		"Age: " + strconv.Itoa(s.Age)

	requestBody, _ := json.Marshal(map[string]string{
		"model":  "llama3",
		"prompt": prompt,
	})

	ollamaURL := "https://0654-2405-201-c054-88da-5c4b-5d55-1c8c-c91e.ngrok-free.app/api/generate"
	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Error calling Ollama", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	summary, _ := result["response"].(string)
	json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}


