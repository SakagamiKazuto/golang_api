package test

import (
	"os"
	"testing"

	"github.com/SakagamiKazuto/golang_api/db"
)

/*
このテストをエントリーポイントとしてすべてのテストが呼ばれる仕組み
処理は以下のように行われる。
1. テストDB接続
2. サンプルデータの挿入
3. テストの実行
 */
func TestMain(m *testing.M) {
	connectDB := db.ConnectTestDB()
	db.DeleteTestData(connectDB)
	db.InsertTestData(connectDB)
	code := m.Run()

	os.Exit(code)
}

