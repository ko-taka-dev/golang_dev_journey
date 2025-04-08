package gui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/client"
	"github.com/ko-taka-dev/golang_dev_journey/todo/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func StartGUI(apiBaseURL string) {
	a := app.New()
	w := a.NewWindow("TODO アプリ")

	todoClient := client.NewTodoClient(apiBaseURL)
	var todos []domain.Todo
	currentFilter := "all"

	var todoList *widget.List
	// タスクのリフレッシュ
	refreshTodos := func() {
		t, err := todoClient.GetTodos()
		if err != nil {
			dialog.ShowError(fmt.Errorf("TODOの取得に失敗しました: %v", err), w)
			return
		}
		todos = t
		todoList.Refresh()
	}

	// フィルタリングされたタスクの取得
	filteredTodos := func() []domain.Todo {
		switch currentFilter {
		case "done":
			var done []domain.Todo
			for _, t := range todos {
				if t.Done {
					done = append(done, t)
				}
			}
			return done
		case "undone":
			var undone []domain.Todo
			for _, t := range todos {
				if !t.Done {
					undone = append(undone, t)
				}
			}
			return undone
		default:
			return todos
		}
	}

	// タスクリストを表示
	todoList = widget.NewList(
		func() int {
			return len(filteredTodos())
		},
		func() fyne.CanvasObject {
			completeCheck := widget.NewCheck("", nil)  // 完了用チェックボックス
			label := widget.NewLabel("")
			left := container.NewHBox(completeCheck, label)

			deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			row := container.NewBorder(nil, nil, nil, deleteBtn, left)

			bg := canvas.NewRectangle(color.Transparent)
			return container.NewStack(bg, row)
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			row := obj.(*fyne.Container).Objects[1].(*fyne.Container)
		
			left := row.Objects[0].(*fyne.Container)
			deleteBtn := row.Objects[1].(*widget.Button)
		
			completeCheck := left.Objects[0].(*widget.Check)
			label := left.Objects[1].(*widget.Label)
		
			filtered := filteredTodos()
			if i >= len(filtered) {
				return // インデックスが範囲外の場合は何もしない
			}
			
			todo := filtered[i] // i番目のタスクを取得
			label.SetText(todo.Title)

			// 重要: OnChangedハンドラを設定する前にチェック状態を設定
			completeCheck.OnChanged = nil // 一時的にハンドラを無効化
			completeCheck.SetChecked(todo.Done)
		
			// OnChangedの中で参照がずれないようにtodoIDを固定
			todoID := strconv.Itoa(int(todo.ID))
		
			completeCheck.OnChanged = func(done bool) {
				go func(id string, status bool) {
					todoClient.PutTodoCompletionStatus(id, status)
					time.Sleep(200 * time.Millisecond)
					refreshTodos()
				}(todoID, done)
			}
		
			deleteBtn.OnTapped = func() {
				dialog.ShowConfirm("確認", "このタスクを削除しますか？", func(confirmed bool) {
					if confirmed {
						todoClient.DeleteTodoByID(todoID)
						refreshTodos()
					}
				}, w)
			}
		},
	)

	input := widget.NewEntry()
	input.SetPlaceHolder("タスクを入力してください")

	addBtn := widget.NewButton("追加", func() {
		if input.Text != "" {
			_, err := todoClient.CreateTodo(input.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("TODOの追加に失敗しました: %v", err), w)
				return
			}
			input.SetText("")
			refreshTodos() // 画面を更新
		}
	})

	inputLine := container.NewBorder(nil, nil, nil, addBtn, input)

	// フィルターボタン
	filterRadio := widget.NewRadioGroup([]string{"全て", "未完了のみ", "完了のみ"}, func(value string) {
		switch value {
		case "全て":
			currentFilter = "all"
		case "未完了のみ":
			currentFilter = "undone"
		case "完了のみ":
			currentFilter = "done"
		}
		todoList.Refresh()
	})
	filterRadio.Horizontal = true
	filterRadio.Selected = "全て" // 初期状態

	header := container.NewHBox(
		canvas.NewText("Todoアプリ", color.White),
		layout.NewSpacer(),
	)

	headerBG := canvas.NewRectangle(color.Black)
	headerBG.SetMinSize(fyne.NewSize(800, 40))
	headerContent := container.NewStack(headerBG, container.NewPadded(header))

	scroll := container.NewScroll(todoList)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	main := container.NewVBox(
		headerContent,
		inputLine,
		filterRadio,
		scroll,
	)

	refreshTodos()
	w.SetContent(main)
	w.Resize(fyne.NewSize(500, 600))
	w.SetFixedSize(false) // ウィンドウサイズ変更を許可
	w.ShowAndRun()
}