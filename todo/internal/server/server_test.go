package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTodoUseCase は usecase.TodoUseCaseInterface のモック実装です
type MockTodoUseCase struct {
	mock.Mock
}

// インターフェースを実装していることを確認
var _ usecase.TodoUseCaseInterface = (*MockTodoUseCase)(nil)

// GetTodos は全てのTodoを取得するメソッドのモックです
func (m *MockTodoUseCase) GetTodos() []domain.Todo {
	args := m.Called()
	return args.Get(0).([]domain.Todo)
}

// CreateTodo は新しいTodoを作成するメソッドのモックです
func (m *MockTodoUseCase) CreateTodo(title string) domain.Todo {
	args := m.Called(title)
	return args.Get(0).(domain.Todo)
}

// DeleteTodoByID はIDを指定してTodoを削除するメソッドのモックです
func (m *MockTodoUseCase) DeleteTodoByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// CompleteTodoByID はIDを指定してTodoを完了状態にするメソッドのモックです
func (m *MockTodoUseCase) CompleteTodoByID(id string) (domain.Todo, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Todo), args.Error(1)
}

func TestGetTodos(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	server := NewTodoServer(mockUseCase)

	expectedTodos := []domain.Todo{{ID: 1, Title: "Test Todo", Done: false}}
	mockUseCase.On("GetTodos").Return(expectedTodos)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	server.getTodos(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response []domain.Todo
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedTodos, response)
}

func TestCreateTodo(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	server := NewTodoServer(mockUseCase)

	newTodo := domain.Todo{ID: 1, Title: "New Todo", Done: false}
	mockUseCase.On("CreateTodo", "New Todo").Return(newTodo)

	reqBody := bytes.NewBufferString(`{"title": "New Todo"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", reqBody)
	w := httptest.NewRecorder()

	server.createTodo(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response domain.Todo
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, newTodo, response)
}

func TestCreateTodo_InvalidJSON(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	server := NewTodoServer(mockUseCase)

	reqBody := bytes.NewBufferString(`invalid json`)
	req := httptest.NewRequest(http.MethodPost, "/todos", reqBody)
	w := httptest.NewRecorder()

	server.createTodo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTodo(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	server := NewTodoServer(mockUseCase)

	mockUseCase.On("DeleteTodoByID", "1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	server.deleteTodo(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteTodo_NotFound(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	server := NewTodoServer(mockUseCase)

	mockUseCase.On("DeleteTodoByID", "999").Return(errors.New("todo not found"))

	req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	server.deleteTodo(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCompleteTodo(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	server := NewTodoServer(mockUseCase)

	completedTodo := domain.Todo{ID: 1, Title: "Test Todo", Done: true}
	mockUseCase.On("CompleteTodoByID", "1").Return(completedTodo, nil)

	req := httptest.NewRequest(http.MethodPut, "/todos/1/done", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	server.completeTodo(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response domain.Todo
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, completedTodo, response)
}

func TestCompleteTodo_NotFound(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	server := NewTodoServer(mockUseCase)

	mockUseCase.On("CompleteTodoByID", "999").Return(domain.Todo{}, errors.New("todo not found"))

	req := httptest.NewRequest(http.MethodPut, "/todos/999/done", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	server.completeTodo(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}