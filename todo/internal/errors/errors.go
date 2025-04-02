package errors

import (
    "fmt"
)

// エラー種別を定義
const (
    NotFound      = "NOT_FOUND"
    InvalidInput  = "INVALID_INPUT"
    InternalError = "INTERNAL_ERROR"
)

// AppError はアプリケーション固有のエラー情報を保持
type AppError struct {
    Type    string // エラーの種別
    Message string // エラーメッセージ
    Err     error  // 元のエラー（オプション）
}

// Error はエラーメッセージを返す
func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap は元のエラーを返す
func (e *AppError) Unwrap() error {
    return e.Err
}

// NewNotFoundError は「リソースが見つからない」エラーを作成
func NewNotFoundError(message string, err ...error) *AppError {
    var originalErr error
    if len(err) > 0 {
        originalErr = err[0]
    }
    return &AppError{
        Type:    NotFound,
        Message: message,
        Err:     originalErr,
    }
}

// NewInvalidInputError は「無効な入力」エラーを作成
func NewInvalidInputError(message string, err ...error) *AppError {
    var originalErr error
    if len(err) > 0 {
        originalErr = err[0]
    }
    return &AppError{
        Type:    InvalidInput,
        Message: message,
        Err:     originalErr,
    }
}

// NewInternalError は「内部エラー」を作成
func NewInternalError(message string, err ...error) *AppError {
    var originalErr error
    if len(err) > 0 {
        originalErr = err[0]
    }
    return &AppError{
        Type:    InternalError,
        Message: message,
        Err:     originalErr,
    }
}

// IsNotFound はエラーが「リソースが見つからない」エラーかどうかを判定
func IsNotFound(err error) bool {
    var appErr *AppError
    if err == nil {
        return false
    }
    if as, ok := err.(*AppError); ok {
        appErr = as
        return appErr.Type == NotFound
    }
    return false
}

// IsInvalidInput はエラーが「無効な入力」エラーかどうかを判定
func IsInvalidInput(err error) bool {
    var appErr *AppError
    if err == nil {
        return false
    }
    if as, ok := err.(*AppError); ok {
        appErr = as
        return appErr.Type == InvalidInput
    }
    return false
}

// IsInternalError はエラーが「内部エラー」かどうかを判定
func IsInternalError(err error) bool {
    var appErr *AppError
    if err == nil {
        return false
    }
    if as, ok := err.(*AppError); ok {
        appErr = as
        return appErr.Type == InternalError
    }
    return false
}