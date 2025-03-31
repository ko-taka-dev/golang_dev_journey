package gui

import (
    "fmt"
    "strconv"
    "github.com/ko-taka-dev/golang_dev_journey/todo/internal/usecase"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func StartGUI(usecase *usecase.TodoUseCase) {
    a := app.New()
    w := a.NewWindow("TODO アプリ")
    w.Resize(fyne.NewSize(400, 600))

    var selectedID *uint // 選択中のTODOのID（NULL許容）

    todoList := widget.NewList(
        func() int {
            return len(usecase.GetTodos())
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(i widget.ListItemID, obj fyne.CanvasObject) {
            todos := usecase.GetTodos()
            obj.(*widget.Label).SetText(fmt.Sprintf("%d: %s [%t]", todos[i].ID, todos[i].Title, todos[i].Done))
        },
    )

    // アイテムが選択されたときの処理
    todoList.OnSelected = func(id widget.ListItemID) {
        todos := usecase.GetTodos()
        if id >= 0 && id < len(todos) {
            selectedID = &todos[id].ID
        }
    }

    input := widget.NewEntry()
    input.SetPlaceHolder("新しい TODO を入力")

    addButton := widget.NewButton("追加", func() {
        if input.Text != "" {
            usecase.CreateTodo(input.Text)
            input.SetText("")
            todoList.Refresh()
        }
    })

    deleteButton := widget.NewButton("削除", func() {
        if selectedID != nil {
            idStr := strconv.Itoa(int(*selectedID))
            usecase.DeleteTodoByID(idStr)
            selectedID = nil // 削除後はリセット
            todoList.Refresh()
        }
    })

    completeButton := widget.NewButton("完了", func() {
        if selectedID != nil {
            idStr := strconv.Itoa(int(*selectedID))
            usecase.CompleteTodoByID(idStr)
            selectedID = nil // 完了後はリセット
            todoList.Refresh()
        }
    })

    w.SetContent(container.NewVBox(
        input,
        addButton,
        todoList,
        container.NewHBox(deleteButton, completeButton),
    ))

    w.ShowAndRun()
}