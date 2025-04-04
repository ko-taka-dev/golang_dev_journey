package repository

import (
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"gorm.io/gorm"
)

// TodoRepository はTodoエンティティのデータアクセスを担当する構造体
type TodoRepository struct {
    db *gorm.DB
}

// TodoRepositoryInterface はTodoRepositoryのインターフェース
type TodoRepositoryInterface interface {
	FindAll() ([]domain.Todo, error)
	FindByID(id string) (*domain.Todo, error)
	Create(todo *domain.Todo) error
	Update(todo *domain.Todo) error
	Delete(todo *domain.Todo) error
}

// NewTodoRepository はTodoRepositoryのコンストラクタ
func NewTodoRepository(db *gorm.DB) TodoRepositoryInterface {
    return &TodoRepository{db: db}
}

// FindAll はすべてのTodoを取得するメソッド
func (r *TodoRepository) FindAll() ([]domain.Todo, error) {
	var todos []domain.Todo
	result := r.db.Find(&todos)
	return todos, result.Error
}

// FindByID は指定されたIDのTodoを取得するメソッド
func (r *TodoRepository) FindByID(id string) (*domain.Todo, error) {
	var todo domain.Todo
	result := r.db.First(&todo, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // レコードが見つからない場合は特別扱い
		}
		return nil, result.Error
	}
	return &todo, nil
}

// Create は新しいTodoを作成するメソッド
func (r *TodoRepository) Create(todo *domain.Todo) error {
    result := r.db.Create(todo)
    return result.Error
}

// Update は指定されたTodoを更新するメソッド
func (r *TodoRepository) Update(todo *domain.Todo) error {
    result := r.db.Save(todo)
    return result.Error
}

// Delete は指定されたTodoを削除するメソッド
func (r *TodoRepository) Delete(todo *domain.Todo) error {
    result := r.db.Delete(todo)
    return result.Error
}
