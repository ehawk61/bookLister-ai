package model

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func NewID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (b Book) Validate() string {
	var problems []string
	if strings.TrimSpace(b.Title) == "" {
		problems = append(problems, "title is required")
	}
	if strings.TrimSpace(b.Author) == "" {
		problems = append(problems, "author is required")
	}
	return strings.Join(problems, "; ")
}
