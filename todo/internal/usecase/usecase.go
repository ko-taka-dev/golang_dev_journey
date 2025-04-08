package usecase

import (
	"fmt"
	"strings"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/repository"
	"github.com/ko-taka-dev/golang_dev_journey/todo/pkg/errors"
)

// TodoUseCaseInterface はTodoのビジネスロジックを定義するインターフェース
type TodoUseCaseInterface interface {
    GetTodos() ([]domain.Todo, error)
    CreateTodo(title string) (domain.Todo, error)
    UpdateTodo(id string, done bool) (domain.Todo, error)
    DeleteTodoByID(id string) error
}

// TodoUseCase は TodoUseCaseInterface を実装する構造体
type TodoUseCase struct {
    repo repository.TodoRepositoryInterface // インターフェースを使う
}

// NewTodoUseCase は新しいTodoUseCaseインスタンスを作成する関数
func NewTodoUseCase(repo repository.TodoRepositoryInterface) TodoUseCaseInterface {
    return &TodoUseCase{repo: repo}
}

// GetTodos はすべてのTODOを取得するメソッド
func (uc *TodoUseCase) GetTodos() ([]domain.Todo, error) {
	todos, err := uc.repo.FindAll()
	if err != nil {
		return nil, errors.NewInternalError("Todoの取得に失敗しました", err)
	}
    return todos, nil
}

// CreateTodo は新しいTODOを作成するメソッド
func (uc *TodoUseCase) CreateTodo(title string) (domain.Todo, error) {
    // タイトルの検証
    if title == "" {
        return domain.Todo{}, errors.NewInvalidInputError("タイトルは必須です")
    }

    // タイトルの長さ制限
    if len(title) > 100 {
        return domain.Todo{}, errors.NewInvalidInputError("タイトルは100文字以内にしてください")
    }
    
    // タイトルのトリミング（前後の空白を削除）
    title = strings.TrimSpace(title)
    // トリミング後に空になった場合もエラー
    if title == "" {
        return domain.Todo{}, errors.NewInvalidInputError("タイトルに有効な文字を入力してください")
    }

    todo := domain.Todo{Title: title, Done: false}
    if err := uc.repo.Create(&todo); err != nil {
        return domain.Todo{}, errors.NewInternalError("Todoの作成に失敗しました", err)
    }
    return todo, nil
}

// UpdateTodo は指定されたIDのTODOを更新するメソッド
func (uc *TodoUseCase) UpdateTodo(id string, done bool) (domain.Todo, error) {
	todo, err := uc.repo.FindByID(id)
	if err != nil {
		return domain.Todo{}, errors.NewInternalError(fmt.Sprintf("ID %s のTodoの検索に失敗しました", id), err)
	}
	
	if todo == nil {
		return domain.Todo{}, errors.NewNotFoundError(fmt.Sprintf("ID %s のTodoが見つかりません", id))
	}

    todo.Done = done
    if err := uc.repo.Update(todo); err != nil {
		return domain.Todo{}, errors.NewInternalError(fmt.Sprintf("ID %s のTodoの更新に失敗しました", id), err)
    }
    return *todo, nil
}

// DeleteTodoByID は指定されたIDのTODOを削除するメソッド
func (uc *TodoUseCase) DeleteTodoByID(id string) error {
	todo, err := uc.repo.FindByID(id)
	if err != nil {
		return errors.NewInternalError(fmt.Sprintf("ID %s のTodoの検索に失敗しました", id), err)
	}
	
	if todo == nil {
		return errors.NewNotFoundError(fmt.Sprintf("ID %s のTodoが見つかりません", id))
	}
	
	if err := uc.repo.Delete(todo); err != nil {
		return errors.NewInternalError(fmt.Sprintf("ID %s のTodoの削除に失敗しました", id), err)
	}
    return nil
}