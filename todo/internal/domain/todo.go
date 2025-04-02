package domain

// Todo はタスク管理のための基本的なデータ構造
type Todo struct {
    ID    uint   `gorm:"primaryKey"` // タスクの一意識別子
    Title string `json:"title"`  // タスクのタイトル
    Done  bool   `json:"done"`   // タスクの完了状態（true: 完了、false: 未完了）
}
