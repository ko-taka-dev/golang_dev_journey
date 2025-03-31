package calculator

import "fmt"

// 計算処理
func Calculate(a float64, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("エラー: 0で割ることはできません")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("エラー: '%s' は無効な演算子です", operator)
	}
}
