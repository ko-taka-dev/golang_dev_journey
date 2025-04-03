package utils

import (
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

// タイトルの最大文字数
const MaxTitleLength = 100

// SQLインジェクションを防ぐための禁止ワード
var dangerousSQLPatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)\b(INSERT\s+INTO|DELETE\s+FROM|UPDATE\s+\w+\s+SET|DROP\s+TABLE|ALTER\s+TABLE)\b`),
}

// XSS攻撃を防ぐための禁止パターン
var dangerousXSSPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)<script.*?>.*?</script>`), // <script>タグ全体をブロック
	regexp.MustCompile(`(?i)javascript:`),             // javascript:スキーム
	regexp.MustCompile(`(?i)on\w+=".*?"`),             // onclick, onerror などのイベントハンドラ
}

// タイトルのバリデーションを行う
func ValidateTitle(title string) (bool, string) {
	// 空チェック
	if strings.TrimSpace(title) == "" {
		return false, "タイトルは空にできません"
	}

	// 長さチェック
	if utf8.RuneCountInString(title) > MaxTitleLength {
		return false, "タイトルが長すぎます"
	}

	// SQLインジェクションの禁止ワードチェック
	for _, pattern := range dangerousSQLPatterns {
		if pattern.MatchString(title) {
			return false, "タイトルに不正なSQL構文が含まれています"
		}
	}

	// XSSの禁止ワードチェック
	for _, pattern := range dangerousXSSPatterns {
		if pattern.MatchString(title) {
			return false, "タイトルに不正なスクリプトが含まれています"
		}
	}
	return true, ""
}

// HTMLエスケープ（XSS対策）
func SanitizeInput(input string) string {
    return html.EscapeString(input)
}
