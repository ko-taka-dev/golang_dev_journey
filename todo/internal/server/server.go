package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/usecase"
	"github.com/ko-taka-dev/golang_dev_journey/todo/pkg/errors"
	"github.com/ko-taka-dev/golang_dev_journey/todo/pkg/logger"
)

// TodoServer はHTTPリクエストを処理するサーバー
type TodoServer struct {
	router  *mux.Router
	useCase usecase.TodoUseCaseInterface // インターフェースを使用
    logger  *logger.Logger   
}

// NewTodoServer は新しいTodoServerインスタンスを作成する
func NewTodoServer(useCase usecase.TodoUseCaseInterface) *TodoServer {
	s := &TodoServer{
		router:  mux.NewRouter(),
		useCase: useCase,
        logger:  logger.GetLogger(),
	}
	s.routes()
	return s
}

// routes はサーバーのルーティングを設定する
func (s *TodoServer) routes() {
	s.router.HandleFunc("/todos", s.getTodos).Methods("GET")
	s.router.HandleFunc("/todos", s.createTodo).Methods("POST")
	s.router.HandleFunc("/todos/{id}", s.deleteTodo).Methods("DELETE")
	s.router.HandleFunc("/todos/{id}/done", s.completeTodo).Methods("PUT")
}

// Start はサーバーを指定されたアドレスで起動する
func (s *TodoServer) Start(addr string) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s)
}

// ServeHTTP はHTTPリクエストを処理する
func (s *TodoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// getTodos はすべてのTODOを取得する
func (s *TodoServer) getTodos(w http.ResponseWriter, r *http.Request) {
    s.logger.Info("GET /todos リクエストを受信しました")
    todos, err := s.useCase.GetTodos()
    if err != nil {
        s.logger.Errorf("Todoの取得中にエラーが発生しました: %v", err)
        http.Error(w, "Todoの取得中にエラーが発生しました", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    if err := json.NewEncoder(w).Encode(todos); err != nil {
        s.logger.Errorf("Todoのエンコード中にエラーが発生しました: %v", err)
        http.Error(w, "Todoのエンコード中にエラーが発生しました", http.StatusInternalServerError)
        return
    }
    s.logger.Infof("%d 件のTodoを返却しました", len(todos))
}

// createTodo は新しいTODOを作成する
func (s *TodoServer) createTodo(w http.ResponseWriter, r *http.Request) {
    s.logger.Info("POST /todos リクエストを受信しました")
	var req struct {
		Title string `json:"title"`
	}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        s.logger.Errorf("リクエストボディの解析に失敗しました: %v", err)
        http.Error(w, "リクエストボディの解析に失敗しました", http.StatusBadRequest)
        return
    }
    s.logger.Debugf("リクエスト内容: title=%s", req.Title)

    if req.Title == "" {
        s.logger.Error("タイトルは必須です")
        http.Error(w, "タイトルは必須です", http.StatusBadRequest)
        return
    }

    todo, err := s.useCase.CreateTodo(req.Title)
    if err != nil {
        if errors.IsInvalidInput(err) {
            s.logger.Errorf("無効な入力です: %v", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        s.logger.Errorf("Todoの作成中にエラーが発生しました: %v", err)
        http.Error(w, "Todoの作成中にエラーが発生しました", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    
    if err := json.NewEncoder(w).Encode(todo); err != nil {
        s.logger.Errorf("Todoのエンコード中にエラーが発生しました: %v", err)
        http.Error(w, "Todoのエンコード中にエラーが発生しました", http.StatusInternalServerError)
        return
    }
    s.logger.Infof("新しいTodoを作成しました: id=%d, title=%s", todo.ID, todo.Title)
}

// deleteTodo は指定されたTODOを削除する
func (s *TodoServer) deleteTodo(w http.ResponseWriter, r *http.Request) {
    s.logger.Info("DELETE /todos/{id} リクエストを受信しました")
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
        s.logger.Error("IDは必須です")
        http.Error(w, "IDは必須です", http.StatusBadRequest)
        return
    }
    
    err := s.useCase.DeleteTodoByID(id)
    if err != nil {
        if errors.IsNotFound(err) {
            s.logger.Errorf("指定されたTodoが見つかりません: %v", err)
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        s.logger.Errorf("Todoの削除中にエラーが発生しました: %v", err)
        http.Error(w, "Todoの削除中にエラーが発生しました", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
    s.logger.Infof("Todoを削除しました: id=%s", id)
}

// completeTodo は指定されたTODOを完了状態にする
func (s *TodoServer) completeTodo(w http.ResponseWriter, r *http.Request) {
    s.logger.Info("PUT /todos/{id}/done リクエストを受信しました")
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
        s.logger.Error("IDは必須です")
        http.Error(w, "IDは必須です", http.StatusBadRequest)
        return
    }
    
    todo, err := s.useCase.CompleteTodoByID(id)
    if err != nil {
        if errors.IsNotFound(err) {
            s.logger.Errorf("指定されたTodoが見つかりません: %v", err)
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        s.logger.Errorf("Todoの更新中にエラーが発生しました: %v", err)
        http.Error(w, "Todoの更新中にエラーが発生しました", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    if err := json.NewEncoder(w).Encode(todo); err != nil {
        s.logger.Errorf("Todoのエンコード中にエラーが発生しました: %v", err)
        http.Error(w, "Todoのエンコード中にエラーが発生しました", http.StatusInternalServerError)
        return
    }
    s.logger.Infof("Todoを完了状態にしました: id=%s, title=%s", id, todo.Title)
}
