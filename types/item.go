package types

import "time"

type Item struct {
	ID        int       `json:"id"`
	TodoID    int       `json:"todo_id"`
	Name      string    `json:"name"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	Todo      Todo
}

func NewItem(todo_id int, name string, done bool) (*Item, error) {
	return &Item{
		TodoID:    todo_id,
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Done:      done,
	}, nil
}
