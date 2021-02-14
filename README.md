# データ構成
このAPIには「Users」と「Bosyus」の2つのテーブルがあり
UsersがBosyusがhas manyの関係で結ばれています。

# 実装したAPI
Usersに関してはSignup(create)とLogin →[ソースコード](https://github.com/SakagamiKazuto/golang_api/blob/main/handler/auth.go)

Bosyusに関してはCRUDのAPIを作成しました →[ソースコード](https://github.com/SakagamiKazuto/golang_api/blob/main/handler/handler.go)

[](https://golang-api-portfolio.herokuapp.com/swagger/index.html)
