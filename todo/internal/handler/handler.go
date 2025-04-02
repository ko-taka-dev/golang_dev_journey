package handler

import (
	"net/http"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/errors"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/usecase"

	"github.com/gin-gonic/gin"
)

// TodoHandler はTodoのHTTPリクエストを処理する構造体
type TodoHandler struct {
    usecase *usecase.TodoUseCase
}

// NewTodoHandler はTodoHandlerを作成する関数
func NewTodoHandler(usecase *usecase.TodoUseCase) *TodoHandler {
    return &TodoHandler{usecase: usecase}
}

// CreateTodoHandler は新しいTODOを作成するハンドラ
func (h *TodoHandler) CreateTodoHandler(c *gin.Context) {
    var req struct {
        Title string `json:"title"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todo, err := h.usecase.CreateTodo(req.Title)
    if err != nil {
        if errors.IsInvalidInput(err) {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Todoの作成中にエラーが発生しました"})
        return
    }
    c.JSON(http.StatusCreated, todo)
}

// GetTodosHandler はすべてのTODOを取得するハンドラ
func (h *TodoHandler) GetTodosHandler(c *gin.Context) {
    todos, err := h.usecase.GetTodos()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Todoの取得中にエラーが発生しました"})
        return
    }
    c.JSON(http.StatusOK, todos)
}

// CompleteTodoHandler はTODOを完了状態にするハンドラ
func (h *TodoHandler) CompleteTodoHandler(c *gin.Context) {
    id := c.Param("id")
    todo, err := h.usecase.CompleteTodoByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, todo)
}

// DeleteTodoHandler はTODOを削除するハンドラ
func (h *TodoHandler) DeleteTodoHandler(c *gin.Context) {
    id := c.Param("id")
    if err := h.usecase.DeleteTodoByID(id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}
