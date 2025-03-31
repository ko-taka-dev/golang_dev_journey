package repository

import (
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"gorm.io/gorm"
)

type TodoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
    return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *domain.Todo) {
    r.db.Create(todo)
}

func (r *TodoRepository) FindAll() []domain.Todo {
    var todos []domain.Todo
    r.db.Find(&todos)
    return todos
}

func (r *TodoRepository) FindByID(id string) *domain.Todo {
    var todo domain.Todo
    r.db.First(&todo, id)
    if todo.ID == 0 {
        return nil
    }
    return &todo
}

func (r *TodoRepository) Update(todo *domain.Todo) {
    r.db.Save(todo)
}

func (r *TodoRepository) Delete(todo *domain.Todo) {
    r.db.Delete(todo)
}
