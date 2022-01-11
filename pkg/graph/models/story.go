package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DBStoryView struct {
	ID   uuid.UUID       `json:"id"`
	Type string          `json:"type"`
	Spec json.RawMessage `json:"spec"`
}

type DBStory struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Group        string        `json:"group"`
	Created      time.Time     `json:"created"`
	LastModified time.Time     `json:"lastModified"`
	Views        []DBStoryView `json:"views"`
	Draft        bool
}

type GraphStory struct {
	ID           uuid.UUID        `json:"id"`
	Name         string           `json:"name"`
	Group        string           `json:"group"`
	Created      time.Time        `json:"created"`
	LastModified *time.Time       `json:"lastModified"`
	Views        []GraphStoryView `json:"views"`
	Draft        bool
}

type GraphStoryView interface {
	IsStoryView()
}

type StoryViewHeader struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
	Level   int       `json:"level"`
}

func (StoryViewHeader) IsStoryView() {}

type StoryViewMarkdown struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}

func (StoryViewMarkdown) IsStoryView() {}

type StoryViewPlotly struct {
	ID     uuid.UUID                `json:"id"`
	Data   []map[string]interface{} `json:"data"`
	Layout map[string]interface{}   `json:"layout"`
}

func (StoryViewPlotly) IsStoryView() {}
