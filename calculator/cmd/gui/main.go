package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/ko-taka-dev/golang_dev_journey/calculator/internal/calculator"
)

func main() {
	a := app.New()
	w := a.NewWindow("シンプル計算機")
	w.Resize(fyne.NewSize(300, 200))

	entry1 := widget.NewEntry()
	entry1.SetPlaceHolder("数値1")
	entry2 := widget.NewEntry()
	entry2.SetPlaceHolder("数値2")
	operator := widget.NewEntry()
	operator.SetPlaceHolder("演算子 (+, -, *, /)")

	resultLabel := widget.NewLabel("結果: ")

	calculateButton := widget.NewButton("計算", func() {
		a, err1 := strconv.ParseFloat(entry1.Text, 64)
		b, err2 := strconv.ParseFloat(entry2.Text, 64)
		op := operator.Text

		if err1 != nil || err2 != nil {
			resultLabel.SetText("エラー: 数値を入力してください")
			return
		}

		result, err := calculator.Calculate(a, b, op)
		if err != nil {
			resultLabel.SetText(err.Error())
			return
		}

		resultLabel.SetText(fmt.Sprintf("結果: %f", result))
	})

	w.SetContent(container.NewVBox(
		entry1,
		operator,
		entry2,
		calculateButton,
		resultLabel,
	))

	w.ShowAndRun()
}
