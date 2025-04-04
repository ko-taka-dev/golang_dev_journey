package gui

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/client"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// StartGUI GUIを起動
func StartGUI(apiBaseURL string) {
	a := app.New()
	w := a.NewWindow("TODO アプリ")

	todoClient := client.NewTodoClient(apiBaseURL)
	var todos []domain.Todo
	var err error
	var selectedID *uint // 選択中のTODOのID（NULL許容）
	var selectedIndex int = -1   // indexも保持（視覚フィードバック・再選択用）

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
			// 背景色つきセルを構築
			bg := canvas.NewRectangle(color.Transparent) // デフォは透明
			label := widget.NewLabel("")
			return container.NewMax(bg, container.NewHBox(label))
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			// 各セルのデータ更新
			bg := obj.(*fyne.Container).Objects[0].(*canvas.Rectangle)
			label := obj.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
	
			if i < len(todos) {
				status := "未完了"
				if todos[i].Done {
					status = "完了"
				}
				label.SetText(fmt.Sprintf("ID:%d %s [%s]", todos[i].ID, todos[i].Title, status))
	
				// 選択中なら青背景、それ以外は透明
				if i == selectedIndex {
					bg.FillColor = 	color.RGBA{R: 135, G: 206, B: 235, A: 255} // スカイブルー
				} else {
					bg.FillColor = color.Transparent
				}
				bg.Refresh()
			}
		},
	)

	// アイテムが選択されたときの処理
	todoList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(todos) {
			selectedIndex = id
			selectedID = &todos[id].ID
			todoList.Refresh() // 選択の視覚更新

			// 一度選択解除して、再度同じタスクを選びやすくする
			go func(id widget.ListItemID) {
				todoList.Unselect(id)
			}(id)
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
			refreshTodos()
			todoList.Refresh()

			// ★削除後に選択インデックスを補正
			if selectedIndex >= len(todos) {
				selectedIndex = len(todos) - 1
			}
			if selectedIndex >= 0 {
				selectedID = &todos[selectedIndex].ID
			} else {
				selectedID = nil
			}
			todoList.Select(selectedIndex)
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

			// ★完了後も選択保持
			if selectedIndex >= 0 && selectedIndex < len(todos) {
				selectedID = &todos[selectedIndex].ID
				todoList.Select(selectedIndex)
			}
		} else {
			dialog.ShowInformation("選択エラー", "完了するTODOを選択してください", w)
		}
	})

	refreshButton := widget.NewButton("更新", func() {
		refreshTodos()
		todoList.Refresh()
	})

	// 高さ固定＋スクロール化
	todoList.Resize(fyne.NewSize(400, 200))
	scroll := container.NewScroll(todoList)
	scroll.SetMinSize(fyne.NewSize(400, 200)) // 高さ固定
	scrollContainer := container.NewStack(scroll) // サイズ反映させる

	// UI構成（ボタンを下に）
	content := container.NewVBox(
		input,
		addButton,
		scrollContainer, // スクロール可能なTODOリスト
		container.NewHBox(deleteButton, completeButton, refreshButton), // 下に表示
	)
	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}