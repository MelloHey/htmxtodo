package types

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Item      Item
}

func NewTodo(name string) (*Todo, error) {
	return &Todo{
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}, nil
}
