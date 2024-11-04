package Markdown_Notetaking_App

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

var store = "./store"

type Note struct {
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *Note) Save() error {
	err := os.MkdirAll(store, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating store: %s\n", err)
	}
	f, _ := os.Create(store + "/" + n.Title + ".json")
	fmt.Println("Created!")

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	n.CreatedAt = time.Now()
	err = encoder.Encode(n)
	if err != nil {
		return fmt.Errorf("error encoding note: %s", err)
	}
	fmt.Println("Saved!")
	return nil
}

func (n *Note) Get(title string) error {
	var l ListNote
	err := l.get()
	if err != nil {
		return err
	}
	for _, i := range l {
		if i.Title == title {
			*n = i
		}
	}
	return nil
}

type ListNote []Note

func (l *ListNote) get() error {
	file, _ := os.ReadDir(store)
	for _, f := range file {
		var temp Note
		fileJson, err := os.ReadFile(store + "/" + f.Name())
		if err != nil {
			return err
		}
		json.Unmarshal(fileJson, &temp)
		*l = append(*l, temp)
	}
	return nil
}

type GrammarApi struct {
	Key       string    `json:"key"`
	Text      string    `json:"text"`
	SessionId uuid.UUID `json:"session_id"`
}
