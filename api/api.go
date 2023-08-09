package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/MelloHey/htmxtodo/store"
	"github.com/MelloHey/htmxtodo/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type APIServer struct {
	listenAddr string
	store      store.Storage
}

type Film struct {
	Title    string
	Director string
}

func NewAPIServer(listenAddr string, store store.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}

}

func (s *APIServer) Run() {
	r := chi.NewRouter()
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger)
	r.HandleFunc("/todos.json", makeHTTPHandleFunc(s.handleTodo))
	r.HandleFunc("/todos", makeHTTPHandleFunc(s.handleTodos))
	//r.HandleFunc("/todos", s.handleGetTodoHtml)
	r.HandleFunc("/todos/{id}", makeHTTPHandleFunc(s.handleTodo))
	r.HandleFunc("/todos/{id}/item", makeHTTPHandleFunc(s.handleItems))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, r)
	//http.ListenAndServe(":3000", r)
}

func (s *APIServer) handleTodo(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		head := "text/html"
		header := r.Header.Get("Accept")

		if strings.Contains(header, head) {
			return s.handleGetTodoHtml(w, r)
		} else {
			return s.handleGetTodoById(w, r)
		}

	}
	if r.Method == "POST" {
		return s.handleCreateItem(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteTodo(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleTodos(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		head := "text/html"
		header := r.Header.Get("Accept")

		if strings.Contains(header, head) {
			return s.handleGetTodosHtml(w, r)
		} else {
			return s.handleGetTodo(w, r)
		}

	}
	if r.Method == "POST" {
		return s.handleCreateTodo(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteTodo(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetTodoHtml(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.ParseFiles("templates/todo/todo.html"))
	todo, _ := s.store.GetTodosByID(id)

	tmpl.Execute(w, todo)
	return nil
}

// handler function #1 - returns the index.html template, with film data
func (s *APIServer) handleGetTodosHtml(w http.ResponseWriter, r *http.Request) error {

	tmpl := template.Must(template.ParseFiles("templates/todo/index.html"))
	todos, _ := s.store.GetTodos()

	tmpl.Execute(w, todos)
	return nil
}

func (s *APIServer) handleGetTodo(w http.ResponseWriter, r *http.Request) error {
	todos, err := s.store.GetTodos()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, todos)
}

func (s *APIServer) handleGetTodoById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	todos, err := s.store.GetTodosByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, todos)
}

func (s *APIServer) handleDeleteTodo(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	s.store.DeleteTodos(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *APIServer) handleCreateTodo(w http.ResponseWriter, r *http.Request) error {
	req := new(types.Todo)

	header := r.Header.Get("Content-Type")
	if header == "application/x-www-form-urlencoded" {
		r.ParseForm()
		req.Name = r.FormValue("name")
		time.Sleep(1 * time.Second)
		todo, err := types.NewTodo(req.Name)
		if err != nil {
			return err
		}
		if err := s.store.CreateTodo(todo); err != nil {
			return err
		}
		tmpl := template.Must(template.ParseFiles("templates/todo/index.html"))
		tmpl.ExecuteTemplate(w, "todo-list-element", req)
		return nil
	}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	todo, err := types.NewTodo(req.Name)
	if err != nil {
		return err
	}
	if err := s.store.CreateTodo(todo); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, todo)
}

func (s *APIServer) TesthandleCreateTodoHtml(w http.ResponseWriter, r *http.Request) error {
	req := new(types.Todo)
	r.ParseForm()
	req.Name = r.FormValue("name")
	time.Sleep(1 * time.Second)
	todo, err := types.NewTodo(req.Name)
	if err != nil {
		return err
	}
	if err := s.store.CreateTodo(todo); err != nil {
		return err
	}
	tmpl := template.Must(template.ParseFiles("templates/todo/index.html"))
	tmpl.ExecuteTemplate(w, "todo-list-element", req)
	return nil
}

func (s *APIServer) handleItems(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}
		head := "text/html"
		header := r.Header.Get("Accept")

		if strings.Contains(header, head) {
			return s.handleGetItemsHtml(w, r)
		} else {
			item, err := s.store.GetItemsByTodoID(id)
			if err != nil {
				return err
			}
			return WriteJSON(w, http.StatusOK, item)
		}

	}

	if r.Method == "POST" {
		return s.handleCreateItem(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetItemsHtml(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.ParseFiles("templates/todo/todo.html"))
	items, _ := s.store.GetItemsByTodoID(id)

	return tmpl.ExecuteTemplate(w, "item-list-element", items)
}

func (s *APIServer) handleCreateItem(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	req := new(types.Item)
	header := r.Header.Get("Content-Type")
	if header == "application/x-www-form-urlencoded" {
		r.ParseForm()
		req.Name = r.FormValue("name")
		req.TodoID = id
		req.Done = false
		time.Sleep(1 * time.Second)
		item, err := types.NewItem(req.TodoID, req.Name, req.Done)
		if err != nil {
			return err
		}
		if err := s.store.CreateItem(item); err != nil {
			return err
		}
		tmpl := template.Must(template.ParseFiles("templates/todo/todo.html"))
		tmpl.ExecuteTemplate(w, "item-list-element", item)
		return nil
	}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	item, err := types.NewItem(req.TodoID, req.Name, req.Done)
	if err != nil {
		return err
	}
	if err := s.store.CreateItem(item); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, item)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteHtml(w http.ResponseWriter, status int, v any) error {
	//w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)

}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
