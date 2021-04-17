package test

import (
	"os"
	"testing"
)

/*
このテストをエントリーポイントとしてすべてのテストが呼ばれる仕組み
処理は以下のように行われる。
1. テストDB接続
2. サンプルデータの挿入
3. テストの実行
 */
func TestMain(m *testing.M) {
	connectDB := NewTestDB()
	DeleteTestData(connectDB)
	InsertTestData(connectDB)
	code := m.Run()

	os.Exit(code)
}

