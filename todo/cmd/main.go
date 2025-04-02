package main

import (
	"flag"
	"log"
	"time"

	"github.com/ko-taka-dev/golang_dev_journey/todo/infrastructure"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/gui"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/repository"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/server"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/usecase"
)

func main() {
	// コマンドライン引数
	apiPort := flag.String("port", "8080", "API server port")
	flag.Parse()

	// データベース初期化
	db := infrastructure.InitDB()

	// リポジトリ、ユースケース、サーバーの初期化
	todoRepo := repository.NewTodoRepository(db)
	todoUseCase := usecase.NewTodoUseCase(todoRepo)
	todoServer := server.NewTodoServer(todoUseCase)

	// サーバーをgoroutineで起動
	go func() {
		log.Printf("Starting server on port %s...", *apiPort)
		if err := todoServer.Start(":" + *apiPort); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	
	// サーバーが起動するまで少し待つ
	time.Sleep(500 * time.Millisecond)

	// APIのベースURL
	apiBaseURL := "http://localhost:" + *apiPort

	// GUIを起動（メインスレッドで実行）
	log.Println("Starting GUI application...")
	gui.StartGUI(apiBaseURL)
}