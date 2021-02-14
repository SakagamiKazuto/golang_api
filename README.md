# このドキュメントについて
このドキュメントはgolang/echoで作成したAPIについて、その仕様やパッケージ構成、起動手順などについて記したものになります。

# 使用技術一覧
1. golang
2. echo
3. gorm
4. jwt-go
5. testing
6. go-sqlmock

# データ構成とAPIについて
このAPIには「Users」と「Bosyus」の2つのテーブルがあり、Users→Bosyusがhas manyの関係で結ばれています。

また、BosyusのAPIにコールするにはUser経由で生成したJWTのTokenをリクエストヘッダーに持たせる必要がある仕様になっています。

Usersに関してはSignup(create)とLogin → [ソースコード](https://github.com/SakagamiKazuto/golang_api/blob/main/handler/auth.go)

Bosyusに関してはCRUDのAPIを作成しました →[ソースコード](https://github.com/SakagamiKazuto/golang_api/blob/main/handler/handler.go)

また基本的なAPIの仕様はこちらでも確認可能です。

<https://golang-api-portfolio.herokuapp.com/swagger/index.html>

# 起動手順
起動確認は以下の手順で行うことが出来ます
```
git clone https://github.com/SakagamiKazuto/golang_api.git
docker-compose up —build db ※ 一度目の起動がデータ初期化の都合でエラーが出る場合がありますが、もう一度docker-compose upを行うと起動できます。
docker-compose up —build api
```
# API動作確認用
またAPIの動作確認は以下コマンドを順番にご利用いただければデータを用意することなくご確認いただけます。
```
1. Create User
curl -X POST -H "Content-Type: application/json" -d '{"Name": "sample1", "Mail":"sample1@gmail.com", "Password": "123"}' localhost:9999/signup

2. Login User
curl -X POST -H "Content-Type: application/json" -d '{"Name": "sample1", "Mail":"sample1@gmail.com", "Password": "123"}' localhost:9999/login

3. CREATE bosyu
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer <JWT-Token>"  -d '{"title": "sample_title", "about": "sample_about", "pref": "愛媛県", "city": "松山市", "level": "player", "user_id": 1}' localhost:9999/api/bosyu/create

4. GET bosyu
curl -X GET -H "Authorization: Bearer <JWT-Token>" http://localhost:9999/api/bosyu/get?user_id=1

5. UPDATE bosyu
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer <JWT-Token>"  -d '{"title": "sample_title", "about": "sample_about", "pref": "北海道", "city": "松山市", "level": "player", "user_id": 1, "id": 1}' localhost:9999/api/bosyu/update
```
# 言語とフレームワークについて
計画段階で今回利用するデータはテーブル2つとAPIが6つという非常にシンプルな構成であったため、

言語には軽量でシンプルにコードを書くことができ、またかねてより一度利用したいと考えていたGolangを採用しました。

フレームワークについては機能性が高いGinや高速処理が可能なGeegoなどとも比較したのですが、比較的シンプルでかつネット上に情報が充実しているという利点からechoを採用しました。

またORMにはRelationが容易であることや比較的記述がシンプルな点からgormを採用しています。

# ディレクトリとパッケージについて
このAPIは以下のような流れでパッケージの参照を行い、逆参照は行わない方針で実装しました。
```
handler → db
↓  ↓
model
```
またディレクトリ構成は以下のコマンドで確認いただけます。

```
tree -d -I 'data'
.
├── common
├── db
│   └── sql
├── docs
├── handler
├── model
└── test
```

# データベースについて
このAPIにおいてデータベースは3つ登場します。

１つ目がherokuにおけるDB、２つ目がローカルにおいて作動するgolang_mysql_dev、そして３つ目がテスト用に用意されたgolang_mysql_testです。

テスト用のデータベースはgolangにおいてテストDBのmock化を調査したところ、handlerなどを含むテストに適切なライブラリが存在しないことが判明したため用意しました。

なおテストDBのcreateはdocker-compose up dbの初回実行時に行われ、その後テーブル内のデータはテスト実行時に初期化されます。

# テストコードについて
基本的にパッケージhandlerおよびmodelの関数に対して、正常、異常系を網羅するように記述しました。

それぞれのテストコードはTest(FuncName)(PackageName)(Normal|Error)といった規則に基づき命名されています。
その上で、以下のように網羅しているパターンをコメントで明示した上でNormalでは正常系、Errorでは異常系のパターンをテストしています。

また一部では[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)を試しに利用を試みましたが、

発行されたクエリの差分を逐一確認する手法は人間にはかなりつらいということが判明したため、あくまでもデータの動きによるテストを実施しています。
```
/*
SignupTests
Handler:
Normal
1. status201
Error
1. Name, Mail, Passwordのいずれかが空欄
2. MailがすでにUsersのテーブルに存在する
*/
```

なおテストの実行結果は以下のコードで確認できます。
```
docker-compose run api go test -v ./test
```

# OpenAPIについて
OpenAPIにはswaggerを採用しました。

またgolangにおいてはswagというコメントに特定の書式を入力することでswagger.ymlを自動記述するライブラリが存在したため、こちらも今回は使用しています。

# herokuデプロイとコンテナ化について
今回は素早くデプロイが出来て、なおかつ他の方のマシンでも即座に動作するものを作りたいと考えていたので、サーバーにはHeroku、コンテナ化にはdocker/docker-composeを採用しました。
