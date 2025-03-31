package handler

import (
	"net/http"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
    usecase *usecase.TodoUseCase
}

func NewTodoHandler(usecase *usecase.TodoUseCase) *TodoHandler {
    return &TodoHandler{usecase: usecase}
}

func (h *TodoHandler) CreateTodoHandler(c *gin.Context) {
    var req struct {
        Title string `json:"title"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todo := h.usecase.CreateTodo(req.Title)
    c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodosHandler(c *gin.Context) {
    todos := h.usecase.GetTodos()
    c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) CompleteTodoHandler(c *gin.Context) {
    id := c.Param("id")
    todo, err := h.usecase.CompleteTodoByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodoHandler(c *gin.Context) {
    id := c.Param("id")
    if err := h.usecase.DeleteTodoByID(id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}
