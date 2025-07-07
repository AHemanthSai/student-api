package main

import (
	"sync"
)

// Student represents a student's profile
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// students stores all student records in memory
var (
	students = make(map[int]Student)
	mu       sync.RWMutex // Ensures thread-safe access to the students map
)
