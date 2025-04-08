package usecase

import (
	"errors"
	"testing"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/repository"
	appErrors "github.com/ko-taka-dev/golang_dev_journey/todo/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTodoRepository struct {
	mock.Mock
}

var _ repository.TodoRepositoryInterface = (*MockTodoRepository)(nil) // インターフェース適合を保証

func (m *MockTodoRepository) FindAll() ([]domain.Todo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) FindByID(id string) (*domain.Todo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) Create(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Update(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func TestGetTodos(t *testing.T) {
	// 様々なテストケースを実行
	testCases := []struct {
		name          string
		mockBehavior  func(*MockTodoRepository)
		expectedTodos []domain.Todo
		expectedError error
	}{
		{
			name: "正常系: すべてのTodoを取得",
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("FindAll").Return([]domain.Todo{{ID: 1, Title: "Todo 1", Done: false}, {ID: 2, Title: "Todo 2", Done: true}}, nil)
			},
			expectedTodos: []domain.Todo{{ID: 1, Title: "Todo 1", Done: false}, {ID: 2, Title: "Todo 2", Done: true}},
			expectedError: nil,
		},
		{
			name: "異常系: エラーが発生",
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("FindAll").Return(nil, errors.New("データベースエラー"))
			},
			expectedTodos: nil,
			expectedError: appErrors.NewInternalError("Todoの取得に失敗しました", errors.New("データベースエラー")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockTodoRepository)
			tc.mockBehavior(mockRepo)

			uc := NewTodoUseCase(mockRepo)

			todos, err := uc.GetTodos()
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTodos, todos)
			}

			// モックの検証
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateTodo(t *testing.T) {
	// テストケースの定義
	testCases := []struct {
		name          string
		inputTitle    string
		mockBehavior  func(*MockTodoRepository)
		expectedTodo  domain.Todo
		expectedError error
	}{
		{
			name:       "正常系: 有効なタイトルでTodoが作成される",
			inputTitle: "買い物に行く",
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("Create", mock.MatchedBy(func(todo *domain.Todo) bool {
				return todo.Title == "買い物に行く" && !todo.Done
			})).Return(nil)},
			expectedTodo: domain.Todo{
				Title: "買い物に行く",
				Done:  false,
			},
			expectedError: nil,
		},
		{
			name:       "異常系: 空のタイトル",
			inputTitle: "",
			mockBehavior: func(repo *MockTodoRepository) {
				// Create は呼ばれない想定
			},
			expectedTodo:  domain.Todo{},
			expectedError: appErrors.NewInvalidInputError("タイトルは必須です"),
		},
		{
			name:       "異常系: 空白のみのタイトル",
			inputTitle: "   ",
			mockBehavior: func(repo *MockTodoRepository) {
				// Create は呼ばれない想定
			},
			expectedTodo:  domain.Todo{},
			expectedError: appErrors.NewInvalidInputError("タイトルに有効な文字を入力してください"),
		},
		{
			name:       "異常系: タイトルが長すぎる",
			inputTitle: string(make([]rune, 101)),
			mockBehavior: func(repo *MockTodoRepository) {
				// Create は呼ばれない想定
			},
			expectedTodo:  domain.Todo{},
			expectedError: appErrors.NewInvalidInputError("タイトルは100文字以内にしてください"),
		},
		{
			name:       "異常系: 空白のみのタイトル",
			inputTitle: "   ",
			mockBehavior: func(repo *MockTodoRepository) {
				// Create は呼ばれない想定
			},
			expectedTodo:  domain.Todo{},
			expectedError: appErrors.NewInvalidInputError("タイトルに有効な文字を入力してください"),
		},
		{
			name:       "異常系: リポジトリエラー",
			inputTitle: "有効なタイトル",
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("Create", mock.Anything).Return(errors.New("データベースエラー"))
			},
			expectedTodo:  domain.Todo{},
			expectedError: appErrors.NewInternalError("Todoの作成に失敗しました", errors.New("データベースエラー")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックリポジトリのセットアップ
			mockRepo := new(MockTodoRepository)
			tc.mockBehavior(mockRepo)

			// テスト対象のUseCaseを作成
			uc := NewTodoUseCase(mockRepo)

			// テスト実行
			todo, err := uc.CreateTodo(tc.inputTitle)

			// アサーション
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTodo.Title, todo.Title)
				assert.Equal(t, tc.expectedTodo.Done, todo.Done)
			}

			// モックの検証
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateTodoByID(t *testing.T) {
	// 様々なテストケースを実行
	testCases := []struct {
		name     string
		id       string
		done bool
		mockBehavior  func(*MockTodoRepository)
		expectedTodo  *domain.Todo
		expectedError error
	}{
		{
			name:     "正常系: 存在するIDでTodoを完了",
			id:       "1",
			done: true,
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("FindByID", "1").Return(&domain.Todo{ID: 1, Title: "Todo 1", Done: false}, nil)
				repo.On("Update", mock.MatchedBy(func(todo *domain.Todo) bool {
					return todo.ID == 1 && todo.Done == true
				})).Return(nil)
			},
			expectedTodo:  &domain.Todo{ID: 1, Title: "Todo 1", Done: true},
			expectedError: nil,
		},
		{
			name:     "異常系: 存在しないIDでTodoを完了",
			id:       "999",
			done: false,
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("FindByID", "999").Return(nil, errors.New("not found"))
			},
			expectedTodo:  nil,
			expectedError: appErrors.NewInternalError("ID 999 のTodoの検索に失敗しました", errors.New("not found")),
		},
		{
			name:     "異常系: リポジトリエラー",
			id:       "1",
			done: false,
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("FindByID", "1").Return(&domain.Todo{ID: 1, Title: "Todo 1", Done: false}, nil)
				repo.On("Update", mock.MatchedBy(func(todo *domain.Todo) bool {
					return todo.ID == 1 && todo.Done == false
				})).Return(errors.New("database error"))
			},
			expectedTodo:  nil,
			expectedError: appErrors.NewInternalError("ID 1 のTodoの更新に失敗しました", errors.New("database error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockTodoRepository)
			tc.mockBehavior(mockRepo)

			uc := NewTodoUseCase(mockRepo)

			todo, err := uc.UpdateTodo(tc.id, tc.done)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.True(t, todo.Done)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteTodoByID(t *testing.T) {
	// 様々なテストケースを実行	
	testCases := []struct {
		name     string
		id       string
		mockBehavior  func(*MockTodoRepository)
		expectedError error
	}{
		{
			name:     "正常系: 存在するIDでTodoを削除",
			id:       "1",
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("FindByID", "1").Return(&domain.Todo{ID: 1, Title: "Todo 1", Done: false}, nil)
				repo.On("Delete", mock.MatchedBy(func(todo *domain.Todo) bool {
					return todo.ID == 1
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "異常系: 存在しないIDでTodoを削除",
			id:       "999",
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("FindByID", "999").Return(nil, errors.New("not found"))
			},
			expectedError: appErrors.NewInternalError("ID 999 のTodoの検索に失敗しました", errors.New("not found")),
		},
		{
			name:     "異常系: リポジトリエラー",
			id:       "1",
			mockBehavior: func(repo *MockTodoRepository) {
				repo.On("FindByID", "1").Return(&domain.Todo{ID: 1, Title: "Todo 1", Done: false}, nil)
				repo.On("Delete", mock.MatchedBy(func(todo *domain.Todo) bool {
					return todo.ID == 1
				})).Return(errors.New("database error"))
			},
			expectedError: appErrors.NewInternalError("ID 1 のTodoの削除に失敗しました", errors.New("database error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockTodoRepository)
			tc.mockBehavior(mockRepo)

			uc := NewTodoUseCase(mockRepo)

			err := uc.DeleteTodoByID(tc.id)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}