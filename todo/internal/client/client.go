package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
)

type TodoClient struct {
	baseURL string
}

// NewTodoClient はTodoClientを作成
func NewTodoClient(baseURL string) *TodoClient {
	return &TodoClient{
		baseURL: baseURL,
	}
}

// GetTodos APIからすべてのTODOを取得
func (c *TodoClient) GetTodos() ([]domain.Todo, error) {
	resp, err := http.Get(c.baseURL + "/todos")
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}
	defer resp.Body.Close()

	var todos []domain.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, fmt.Errorf("failed to decode todos: %w", err)
	}
	return todos, nil
}

// CreateTodo 新しいTODOをAPIを通じて作成
func (c *TodoClient) CreateTodo(title string) (*domain.Todo, error) {
	todo := domain.Todo{Title: title}
	jsonData, err := json.Marshal(todo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal todo: %w", err)
	}

	resp, err := http.Post(c.baseURL+"/todos", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}
	defer resp.Body.Close()

	var createdTodo domain.Todo
	if err := json.NewDecoder(resp.Body).Decode(&createdTodo); err != nil {
		return nil, fmt.Errorf("failed to decode created todo: %w", err)
	}
	return &createdTodo, nil
}

// DeleteTodoByID 指定IDのTODOをAPIを通じて削除
func (c *TodoClient) DeleteTodoByID(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/todos/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete todo request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	defer resp.Body.Close()

    // サーバーはStatusNoContent(204)を返すので、それも成功と見なす
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("failed to delete todo: status code %d", resp.StatusCode)
    }
	return nil
}

// CompleteTodoByID 指定IDのTODOをAPIを通じて完了状態に更新
func (c *TodoClient) CompleteTodoByID(id string) (*domain.Todo, error) {
	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/todos/"+id+"/done", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create complete todo request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to complete todo: %w", err)
	}
	defer resp.Body.Close()

	var updatedTodo domain.Todo
	if err := json.NewDecoder(resp.Body).Decode(&updatedTodo); err != nil {
		return nil, fmt.Errorf("failed to decode updated todo: %w", err)
	}
	return &updatedTodo, nil
}