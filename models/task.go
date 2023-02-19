package models

import (
	"errors"

	"github.com/google/uuid"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var (
	ErrMissingTitle       = errors.New("missing title")
	ErrMissingDescription = errors.New("missing description")
)

func NewTask(title, description string) (*Task, error) {
	if title == "" {
		return nil, ErrMissingTitle
	}

	if description == "" {
		return nil, ErrMissingDescription
	}

	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
	}, nil
}
