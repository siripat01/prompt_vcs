package models

import "time"

// Prompt represents a single version of a prompt
 type Prompt struct {
     ID          string    `json:"id"`
     Version     int       `json:"version"`
     Name        string    `json:"name"`
    Content     string    `json:"content"`
    Variables   []string  `json:"variables"` // Dynamic placeholders like {name}, {date}
    Tags        []string  `json:"tags"`
    Collection  string    `json:"collection"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Commits     []Commit  `json:"commits,omitempty"`
}

// Commit represents a version control entry
type Commit struct {
    Hash      string    `json:"hash"`
    Message   string    `json:"message"`
    Author    string    `json:"author"`
    Timestamp time.Time `json:"timestamp"`
    Diff      string    `json:"diff,omitempty"` // Computed on request
}

// DiffResult represents a comparison between two versions
type DiffResult struct {
    FromVersion int      `json:"from_version"`
    ToVersion   int      `json:"to_version"`
    Additions   []string `json:"additions"`
    Deletions   []string `json:"deletions"`
    Context     []string `json:"context"`
}

// Collection represents a group of prompts
type Collection struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    PromptCount int       `json:"prompt_count"`
}