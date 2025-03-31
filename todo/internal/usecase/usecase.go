package usecase

import (
	"errors"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/repository"
)

type TodoUseCase struct {
    repo *repository.TodoRepository
}

func NewTodoUseCase(repo *repository.TodoRepository) *TodoUseCase {
    return &TodoUseCase{repo: repo}
}

func (uc *TodoUseCase) CreateTodo(title string) domain.Todo {
    todo := domain.Todo{Title: title, Done: false}
    uc.repo.Create(&todo)
    return todo
}

func (uc *TodoUseCase) GetTodos() []domain.Todo {
    return uc.repo.FindAll()
}

func (uc *TodoUseCase) CompleteTodoByID(id string) (*domain.Todo, error) {
    todo := uc.repo.FindByID(id)
    if todo == nil {
        return nil, errors.New("TODOが見つかりません")
    }
    todo.Done = true
    uc.repo.Update(todo)
    return todo, nil
}

func (uc *TodoUseCase) DeleteTodoByID(id string) error {
    todo := uc.repo.FindByID(id)
    if todo == nil {
        return errors.New("TODOが見つかりません")
    }
    uc.repo.Delete(todo)
    return nil
}
