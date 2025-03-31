package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"calculator/internal/calculator"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("シンプルな計算機 (例: 12 + 5)")
	fmt.Println("終了するには 'exit' を入力")

	for {
		fmt.Print("式を入力: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("終了します。")
			break
		}

		parts := strings.Fields(input)
		if len(parts) != 3 {
			fmt.Println("エラー: '数値 演算子 数値' の形式で入力してください")
			continue
		}

		a, err1 := strconv.ParseFloat(parts[0], 64)
		b, err2 := strconv.ParseFloat(parts[2], 64)
		operator := parts[1]

		if err1 != nil || err2 != nil {
			fmt.Println("エラー: 数値の形式が正しくありません")
			continue
		}

		result, err := calculator.Calculate(a, b, operator)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("結果: %f\n", result)
	}
}
