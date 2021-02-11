package test

import (
	"os"
	"testing"
	"work/db"
)

/*
このテストをエントリーポイントとしてすべてのテストが呼ばれる仕組み
処理は以下のように行われる。
1. テストDB接続
2. サンプルデータの挿入
3. テストの実行
4. データの初期化
 */
func TestMain(m *testing.M) {
	db.ConnectTestDB()
	db.DeleteTestData()
	db.InsertTestData()
	code := m.Run()

	db.DeleteTestData()
	os.Exit(code)
}

