# Todoアプリ (Golang)

## 概要
このプロジェクトは、Golangを用いて開発されたシンプルなTodoアプリです。基本的なCRUD機能を提供し、タスクの管理を行うことができます。

## 主な機能
- タスクの作成
- タスクの取得
- タスクの更新
- タスクの削除

## 環境構築
### 必要なツール
- Go 1.20以上

### セットアップ手順
1. リポジトリをクローンします。
   ```sh
   git clone https://github.com/ko-taka-dev/golang_dev_journey.git
   cd golang_dev_journey/todo
   ```
2. 必要な依存関係をインストールします。
   ```sh
   go mod tidy
   ```
3. アプリを起動します。
   ```sh
   go run cmd/main.go
   ```

## API エンドポイント
| メソッド | エンドポイント | 説明 |
|----------|--------------|------|
| GET | /todos | すべてのタスクを取得 |
| POST | /todos | 新しいタスクを作成 |
| PUT | /todos/{id}/done | タスクを更新 |
| DELETE | /todos/{id} | タスクを削除 |

## テストの実行
```sh
go test ./...
```

