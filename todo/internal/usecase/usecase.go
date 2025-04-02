package usecase

import (
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/errors"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/repository"
)

// TodoUseCaseInterface はTodoのビジネスロジックを定義するインターフェース
type TodoUseCaseInterface interface {
    GetTodos() ([]domain.Todo, error)
    CreateTodo(title string) (domain.Todo, error)
    CompleteTodoByID(id string) (domain.Todo, error)
    DeleteTodoByID(id string) error
}

// TodoUseCase は TodoUseCaseInterface を実装する構造体
type TodoUseCase struct {
    repo *repository.TodoRepository
}

// NewTodoUseCase は新しいTodoUseCaseインスタンスを作成する関数
func NewTodoUseCase(repo *repository.TodoRepository) *TodoUseCase {
    return &TodoUseCase{repo: repo}
}

// CreateTodo は新しいTODOを作成するメソッド
func (uc *TodoUseCase) CreateTodo(title string) (domain.Todo, error) {
    // タイトルの検証
    if title == "" {
        return domain.Todo{}, errors.NewInvalidInputError("タイトルは必須です")
    }
    
    todo := domain.Todo{Title: title, Done: false}
    if err := uc.repo.Create(&todo); err != nil {
        return domain.Todo{}, errors.NewInternalError("Todoの作成に失敗しました", err)
    }
    return todo, nil
}

// GetTodos はすべてのTODOを取得するメソッド
func (uc *TodoUseCase) GetTodos() ([]domain.Todo, error) {
    todos := uc.repo.FindAll()
    return todos, nil
}

// CompleteTodoByID は指定されたIDのTODOを完了状態にするメソッド
func (uc *TodoUseCase) CompleteTodoByID(id string) (domain.Todo, error) {
    todo := uc.repo.FindByID(id)
    if todo == nil {
        return domain.Todo{}, errors.NewNotFoundError("指定されたIDのTODOが見つかりません")
    }
    
    todo.Done = true
    if err := uc.repo.Update(todo); err != nil {
        return domain.Todo{}, errors.NewInternalError("Todoの更新に失敗しました", err)
    }
    return *todo, nil
}

// DeleteTodoByID は指定されたIDのTODOを削除するメソッド
func (uc *TodoUseCase) DeleteTodoByID(id string) error {
    todo := uc.repo.FindByID(id)
    if todo == nil {
        return errors.NewNotFoundError("指定されたIDのTODOが見つかりません")
    }
    
    if err := uc.repo.Delete(todo); err != nil {
        return errors.NewInternalError("Todoの削除に失敗しました", err)
    }
    return nil
}