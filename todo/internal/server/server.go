package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/usecase"
)

type TodoServer struct {
	router  *mux.Router
	useCase *usecase.TodoUseCase
}

func NewTodoServer(useCase *usecase.TodoUseCase) *TodoServer {
	s := &TodoServer{
		router:  mux.NewRouter(),
		useCase: useCase,
	}
	s.routes()
	return s
}

func (s *TodoServer) routes() {
	s.router.HandleFunc("/todos", s.getTodos).Methods("GET")
	s.router.HandleFunc("/todos", s.createTodo).Methods("POST")
	s.router.HandleFunc("/todos/{id}", s.deleteTodo).Methods("DELETE")
	s.router.HandleFunc("/todos/{id}/done", s.completeTodo).Methods("PUT")
}

// Start starts the HTTP server on the given address
func (s *TodoServer) Start(addr string) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s)
}

func (s *TodoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *TodoServer) getTodos(w http.ResponseWriter, r *http.Request) {
	todos := s.useCase.GetTodos()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (s *TodoServer) createTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo := s.useCase.CreateTodo(req.Title)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (s *TodoServer) deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.useCase.DeleteTodoByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *TodoServer) completeTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	todo, err := s.useCase.CompleteTodoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func StartServer(addr string, useCase *usecase.TodoUseCase) error {
	server := NewTodoServer(useCase)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, server)
}