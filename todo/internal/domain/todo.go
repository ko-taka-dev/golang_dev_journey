package domain

type Todo struct {
    ID    uint   `gorm:"primaryKey"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}
