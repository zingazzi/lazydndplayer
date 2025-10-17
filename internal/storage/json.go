// internal/storage/json.go
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// Storage handles character data persistence
type Storage struct {
	FilePath string
}

// NewStorage creates a new storage instance
func NewStorage(filePath string) *Storage {
	return &Storage{
		FilePath: filePath,
	}
}

// Load loads a character from JSON file
func (s *Storage) Load() (*models.Character, error) {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return new character
			return models.NewCharacter(), nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var character models.Character
	err = json.Unmarshal(data, &character)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &character, nil
}

// Save saves a character to JSON file
func (s *Storage) Save(character *models.Character) error {
	// Ensure directory exists
	dir := filepath.Dir(s.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	err = os.WriteFile(s.FilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Export exports character to a specified file
func (s *Storage) Export(character *models.Character, exportPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(exportPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	err = os.WriteFile(exportPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Import imports character from a specified file
func (s *Storage) Import(importPath string) (*models.Character, error) {
	data, err := os.ReadFile(importPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var character models.Character
	err = json.Unmarshal(data, &character)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &character, nil
}

// GetDefaultPath returns the default character file path
func GetDefaultPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "data/character.json"
	}
	return filepath.Join(homeDir, ".lazydndplayer", "character.json")
}
