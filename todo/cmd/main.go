package main

import (
	"github.com/ko-taka-dev/golang_dev_journey/todo/infrastructure"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/gui"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/repository"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/usecase"
)

func main() {
    db := infrastructure.InitDB()
    repo := repository.NewTodoRepository(db)
    usecase := usecase.NewTodoUseCase(repo)
    gui.StartGUI(usecase) // GUI の起動
}
