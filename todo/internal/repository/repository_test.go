package repository

import (
	"testing"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの作成に失敗しました: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("gormの初期化に失敗しました: %v", err)
	}

	return db, mock, func() {
		mockDB.Close()
	}
}

func TestFindAll(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	// モックの設定
	rows := sqlmock.NewRows([]string{"id", "title", "done"}).
		AddRow(1, "Test Todo 1", false).
		AddRow(2, "Test Todo 2", true)

	mock.ExpectQuery("^SELECT (.+) FROM `todos`").WillReturnRows(rows)

	// テスト対象のリポジトリを作成
	repo := NewTodoRepository(db)

	// テスト実行
	todos, err := repo.FindAll()

	// 検証
	assert.NoError(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, uint(1), todos[0].ID)
	assert.Equal(t, "Test Todo 1", todos[0].Title)
	assert.False(t, todos[0].Done)
	assert.Equal(t, uint(2), todos[1].ID)
	assert.Equal(t, "Test Todo 2", todos[1].Title)
	assert.True(t, todos[1].Done)

	// モックの期待通りに呼ばれたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未実行のクエリがあります: %v", err)
	}
}

func TestFindByID(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	
	// モックの設定 - 存在するID
	rows := sqlmock.NewRows([]string{"id", "title", "done"}).
		AddRow(1, "Test Todo", false)

	mock.ExpectQuery("^SELECT \\* FROM `todos` WHERE `todos`.`id` = \\? ORDER BY `todos`.`id` LIMIT \\?").
		WithArgs("1", 1).
		WillReturnRows(rows)

	// モックの設定 - 存在しないID
	mock.ExpectQuery("^SELECT \\* FROM `todos` WHERE `todos`.`id` = \\? ORDER BY `todos`.`id` LIMIT \\?").
		WithArgs("999", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// テスト対象のリポジトリを作成
	repo := NewTodoRepository(db)

	// テスト実行 - 存在するID
	todo, err := repo.FindByID("1")

	// 検証 - 存在するID
	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, uint(1), todo.ID)
	assert.Equal(t, "Test Todo", todo.Title)
	assert.False(t, todo.Done)

	// テスト実行 - 存在しないID
	todo, err = repo.FindByID("999")

	// 検証 - 存在しないID
	assert.NoError(t, err) // エラーではなくnilを返す設計
	assert.Nil(t, todo)

	// モックの期待通りに呼ばれたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未実行のクエリがあります: %v", err)
	}
}

func TestCreate(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	// モックの設定
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `todos`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// テスト対象のリポジトリを作成
	repo := NewTodoRepository(db)

	// テスト実行
	todo := &domain.Todo{Title: "New Todo", Done: false}
	err := repo.Create(todo)

	// 検証
	assert.NoError(t, err)
	assert.Equal(t, uint(1), todo.ID) // GORMがIDを設定する

	// モックの期待通りに呼ばれたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未実行のクエリがあります: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	// モックの設定
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `todos`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// テスト対象のリポジトリを作成
	repo := NewTodoRepository(db)

	// テスト実行
	todo := &domain.Todo{ID: 1, Title: "Updated Todo", Done: true}
	err := repo.Update(todo)

	// 検証
	assert.NoError(t, err)

	// モックの期待通りに呼ばれたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未実行のクエリがあります: %v", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	// モックの設定
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM `todos`").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// テスト対象のリポジトリを作成
	repo := NewTodoRepository(db)

	// テスト実行
	todo := &domain.Todo{ID: 1, Title: "Test Todo", Done: false}
	err := repo.Delete(todo)

	// 検証
	assert.NoError(t, err)

	// モックの期待通りに呼ばれたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未実行のクエリがあります: %v", err)
	}
}