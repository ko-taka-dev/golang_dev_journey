package infrastructure

import (
    "log"
    "github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// InitDB はデータベースを初期化し、マイグレーションを実行する関数
func InitDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("データベース接続失敗:", err)
    }
    db.AutoMigrate(&domain.Todo{})
    return db
}
