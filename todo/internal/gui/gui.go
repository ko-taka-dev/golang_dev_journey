package gui

import (
	"fmt"
	"strconv"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/client"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func StartGUI(apiBaseURL string) {
	a := app.New()
	w := a.NewWindow("TODO アプリ")
	w.Resize(fyne.NewSize(400, 600))

	todoClient := client.NewTodoClient(apiBaseURL)
	var todos []domain.Todo
	var err error
	var selectedID *uint // 選択中のTODOのID（NULL許容）

	// 初期データ取得
	refreshTodos := func() {
		todos, err = todoClient.GetTodos()
		if err != nil {
			dialog.ShowError(fmt.Errorf("TODOの取得に失敗しました: %v", err), w)
		}
	}

	refreshTodos()

    todoList := widget.NewList(
        func() int {
            return len(todos)
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(i widget.ListItemID, obj fyne.CanvasObject) {
            if i < len(todos) {
                // タイトルが表示されるように確認
                status := "未完了"
                if todos[i].Done {
                    status = "完了"
                }
                obj.(*widget.Label).SetText(fmt.Sprintf("ID:%d %s [%s]", todos[i].ID, todos[i].Title, status))
            }
        },
    )

	// アイテムが選択されたときの処理
	todoList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(todos) {
			selectedID = &todos[id].ID
		}
	}

	input := widget.NewEntry()
	input.SetPlaceHolder("新しい TODO を入力")

	addButton := widget.NewButton("追加", func() {
		if input.Text != "" {
			createdTodo, err := todoClient.CreateTodo(input.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("TODOの追加に失敗しました: %v", err), w)
				return
			}
			input.SetText("")
			refreshTodos()

            // 新しく追加したTODOを選択状態にする
            // 新しいTODOは通常リストの最後に追加されるので、そのインデックスを見つける
            for i, todo := range todos {
                if todo.ID == createdTodo.ID {
                    todoList.Select(i)
                    selectedID = &todo.ID // 選択状態を更新
                    break
                }
            }
            todoList.Refresh()
		}
	})

	deleteButton := widget.NewButton("削除", func() {
		if selectedID != nil {
			idStr := strconv.Itoa(int(*selectedID))
			err := todoClient.DeleteTodoByID(idStr)
			if err != nil {
				dialog.ShowError(fmt.Errorf("TODOの削除に失敗しました: %v", err), w)
				return
			}
			selectedID = nil // 削除後はリセット
			refreshTodos()
			todoList.Refresh()
		} else {
			dialog.ShowInformation("選択エラー", "削除するTODOを選択してください", w)
		}
	})

	completeButton := widget.NewButton("完了", func() {
		if selectedID != nil {
			idStr := strconv.Itoa(int(*selectedID))
			_, err := todoClient.CompleteTodoByID(idStr)
			if err != nil {
				dialog.ShowError(fmt.Errorf("TODOの完了処理に失敗しました: %v", err), w)
				return
			}
			refreshTodos()
			todoList.Refresh()
		} else {
			dialog.ShowInformation("選択エラー", "完了するTODOを選択してください", w)
		}
	})

	refreshButton := widget.NewButton("更新", func() {
		refreshTodos()
		todoList.Refresh()
	})

	w.SetContent(container.NewVBox(
		input,
		addButton,
		todoList,
		container.NewHBox(deleteButton, completeButton, refreshButton),
	))

	w.ShowAndRun()
}