package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/errors"
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
func (m *MockTodoUseCase) GetTodos() ([]domain.Todo, error) {
	args := m.Called()
	return args.Get(0).([]domain.Todo), args.Error(1)
}

// CreateTodo は新しいTodoを作成するメソッドのモックです
func (m *MockTodoUseCase) CreateTodo(title string) (domain.Todo, error) {
	args := m.Called(title)
	return args.Get(0).(domain.Todo), args.Error(1)
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
    // テストケース
    testCases := []struct {
        name        string
        todos       []domain.Todo
        err         error
        expectedStatus int
    }{
        {
            name: "正常系",
            todos: []domain.Todo{{ID: 1, Title: "Test Todo", Done: false}},
            err: nil,
            expectedStatus: http.StatusOK,
        },
        {
            name: "エラー発生",
            todos: []domain.Todo{},
            err: errors.NewInternalError("データベースエラー"),
            expectedStatus: http.StatusInternalServerError,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // モックの設定
            mockUseCase := new(MockTodoUseCase)
            mockUseCase.On("GetTodos").Return(tc.todos, tc.err)
            server := NewTodoServer(mockUseCase)

            // リクエスト実行
            req := httptest.NewRequest(http.MethodGet, "/todos", nil)
            w := httptest.NewRecorder()
            server.getTodos(w, req)

            // 検証
            assert.Equal(t, tc.expectedStatus, w.Code)
            
            if tc.err == nil {
                assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
                var response []domain.Todo
                json.Unmarshal(w.Body.Bytes(), &response)
                assert.Equal(t, tc.todos, response)
            }
            
            mockUseCase.AssertExpectations(t)
        })
    }
}

func TestCreateTodo(t *testing.T) {
    // テストケース
    testCases := []struct {
        name        string
        requestBody string
        todo        domain.Todo
        err         error
        expectedStatus int
    }{
        {
            name: "正常系",
            requestBody: `{"title": "New Todo"}`,
            todo: domain.Todo{ID: 1, Title: "New Todo", Done: false},
            err: nil,
            expectedStatus: http.StatusCreated,
        },
        {
            name: "無効なJSON",
            requestBody: `invalid json`,
            todo: domain.Todo{},
            err: nil,
            expectedStatus: http.StatusBadRequest,
        },
        {
            name: "タイトル未入力",
            requestBody: `{"title": ""}`,
            todo: domain.Todo{},
            err: errors.NewInvalidInputError("タイトルは必須です"),
            expectedStatus: http.StatusBadRequest,
        },
        {
            name: "内部エラー",
            requestBody: `{"title": "New Todo"}`,
            todo: domain.Todo{},
            err: errors.NewInternalError("データベースエラー"),
            expectedStatus: http.StatusInternalServerError,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // モックの設定
            mockUseCase := new(MockTodoUseCase)
            
            // 無効なJSONの場合はCreateTodoが呼ばれないことを期待
            if tc.requestBody != "invalid json" {
                mockUseCase.On("CreateTodo", mock.AnythingOfType("string")).Return(tc.todo, tc.err)
            }
            
            server := NewTodoServer(mockUseCase)

            // リクエスト実行
            reqBody := bytes.NewBufferString(tc.requestBody)
            req := httptest.NewRequest(http.MethodPost, "/todos", reqBody)
            w := httptest.NewRecorder()
            server.createTodo(w, req)

            // 検証
            assert.Equal(t, tc.expectedStatus, w.Code)
            
            if tc.err == nil && tc.requestBody != "invalid json" {
                assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
                var response domain.Todo
                json.Unmarshal(w.Body.Bytes(), &response)
                assert.Equal(t, tc.todo, response)
            }
            
            mockUseCase.AssertExpectations(t)
        })
    }
}

func TestDeleteTodo(t *testing.T) {
    // テストケース
    testCases := []struct {
        name        string
        id          string
        err         error
        expectedStatus int
    }{
        {
            name: "正常系",
            id: "1",
            err: nil,
            expectedStatus: http.StatusNoContent,
        },
        {
            name: "存在しないID",
            id: "999",
            err: errors.NewNotFoundError("指定されたIDのTODOが見つかりません"),
            expectedStatus: http.StatusNotFound,
        },
        {
            name: "内部エラー",
            id: "1",
            err: errors.NewInternalError("データベースエラー"),
            expectedStatus: http.StatusInternalServerError,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // モックの設定
            mockUseCase := new(MockTodoUseCase)
            mockUseCase.On("DeleteTodoByID", tc.id).Return(tc.err)
            server := NewTodoServer(mockUseCase)

            // リクエスト実行
            req := httptest.NewRequest(http.MethodDelete, "/todos/"+tc.id, nil)
            req = mux.SetURLVars(req, map[string]string{"id": tc.id})
            w := httptest.NewRecorder()
            server.deleteTodo(w, req)

            // 検証
            assert.Equal(t, tc.expectedStatus, w.Code)
            mockUseCase.AssertExpectations(t)
        })
    }
}

func TestCompleteTodo(t *testing.T) {
    // テストケース
    testCases := []struct {
        name        string
        id          string
        todo        domain.Todo
        err         error
        expectedStatus int
    }{
        {
            name: "正常系",
            id: "1",
            todo: domain.Todo{ID: 1, Title: "Test Todo", Done: true},
            err: nil,
            expectedStatus: http.StatusOK,
        },
        {
            name: "存在しないID",
            id: "999",
            todo: domain.Todo{},
            err: errors.NewNotFoundError("指定されたIDのTODOが見つかりません"),
            expectedStatus: http.StatusNotFound,
        },
        {
            name: "内部エラー",
            id: "1",
            todo: domain.Todo{},
            err: errors.NewInternalError("データベースエラー"),
            expectedStatus: http.StatusInternalServerError,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // モックの設定
            mockUseCase := new(MockTodoUseCase)
            mockUseCase.On("CompleteTodoByID", tc.id).Return(tc.todo, tc.err)
            server := NewTodoServer(mockUseCase)

            // リクエスト実行
            req := httptest.NewRequest(http.MethodPut, "/todos/"+tc.id+"/done", nil)
            req = mux.SetURLVars(req, map[string]string{"id": tc.id})
            w := httptest.NewRecorder()
            server.completeTodo(w, req)

            // 検証
            assert.Equal(t, tc.expectedStatus, w.Code)
            
            if tc.err == nil {
                assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
                var response domain.Todo
                json.Unmarshal(w.Body.Bytes(), &response)
                assert.Equal(t, tc.todo, response)
            }
            
            mockUseCase.AssertExpectations(t)
        })
    }
}