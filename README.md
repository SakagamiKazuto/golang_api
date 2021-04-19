# このドキュメントについて
このドキュメントはgolang/echoで作成したAPIについて、その仕様やパッケージ構成、起動手順などについて記したものになります。

# 使用技術一覧
1. golang
2. echo
3. gorm
4. jwt-go
5. testing
6. swagger
7. docker
8. heroku
9. air/delve

# エディタについて
エディタについてはgolandを推奨します。

他のエディタを活用する場合、 以下の設定を行ってください。

1. コード補完のためのGopath
2. デバッガのためのremote debug

# データ構成とAPIについて
このAPIには「Users」と「Bosyus」の2つのテーブルがあり、Users→Bosyusがhas manyの関係で結ばれています。

また、BosyusのAPIにコールするにはUser経由で生成したJWTのTokenをリクエストヘッダーに持たせる必要がある仕様になっています。

また基本的なAPIの仕様はこちらでも確認可能です。

<https://goapi-staging.herokuapp.com/swagger/index.html>

# 起動手順
## コンテナの起動
コンテナの起動は以下の手順で行うことが出来ます
```
git clone https://github.com/SakagamiKazuto/golang_api.git
docker-compose up —build db
docker-compose up —build api
```
## ホットリロードとデバッガ
またAPIではホットリロードには[air](https://github.com/cosmtrek/air) ,デバッガには[delve](https://github.com/go-delve/delve) を採用しています。

delveではデバッガの起動時にソースコードの実行が行われるので上記の手順でコンテナ起動を行った後にそのままリクエストを発行した場合以下が返ってきてしまいます。
```
curl: (52) Empty reply from server
```
そのため使用するエディタでデバッガを使用できる環境を設定する必要があります。
デバッガの設定については各エディタにてremote debugの設定を行ってください。

#### idea製品
<https://qiita.com/4486/items/d1dad30403348004fc0a#goland%E3%81%8B%E3%82%89%E3%83%87%E3%83%90%E3%83%83%E3%82%B0%E3%81%99%E3%82%8B>

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
curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer <JWT-Token>" -d '{"user_id": 1}' http://localhost:9999/api/bosyu/get

5. UPDATE bosyu
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer <JWT-Token>"  -d '{"title": "sample_title", "about": "sample_about", "pref": "北海道", "city": "松山市", "level": "player", "user_id": 1, "id": 1}' localhost:9999/api/bosyu/update
```
# 言語とフレームワークについて
計画段階で今回利用するデータはテーブル2つとAPIが6つという非常にシンプルな構成であったため、

言語には軽量でシンプルにコードを書くことができ、またかねてより一度利用したいと考えていたGolangを採用しました。

フレームワークについては機能性が高いGinや高速処理が可能なGeegoなどとも比較したのですが、比較的シンプルでかつネット上に情報が充実しているという利点からechoを採用しました。

またORMにはRelationが容易であることや比較的記述がシンプルな点からgormを採用しています。

# モジュール管理について
このAPIではモジュール管理にはgo.modを使用しています。

コード内から参照するライブラリを使用する場合は一般的なgo modの使用方法にて行い、

デバッガなどのコード内から参照しないツールを導入する場合はDockerfile.devに適宜追加してください。

(例)
```
RUN go get -u github.com/kazukousen/gouml/cmd/gouml
RUN go get -u github.com/cosmtrek/air
RUN go get github.com/go-delve/delve/cmd/dlv
```

# アーキテクチャについて
アーキテクチャにはクリーンアーキテクチャを使用しています。
```
※goumlで生成したクラス図の画像※
```

## infra
### waf
echoに依存する部分を記述します。
この使用されるパッケージ(logger, errorなど)はこのディレクトリ以下に配置してください

### dbhandle
gormに依存する部分を記述します。

## interface
### controller
usecaseを使用してwaf層に返還します。

### database
dbhandle層を使った処理を記述します。
dbhandle層はinterfaceを介して参照されます。

## usecase
database層を使ってdomain層とやりとりを行います。
基本的にはここの処理がcontroller層呼び出されます。

## domain
データベースの構造体が定義されます。

ディレクトリ構成
```
tree -d -I "pkg|src|data|tmp"
.
├── docs
├── domain
├── infra
│   ├── dbhandle
│   │   ├── logs
│   │   └── sql
│   └── waf
│       ├── apperror
│       └── logger
├── interface
│   ├── controller
│   └── database
├── test
└── usecasee
```

# 例外処理について
外部起因および内部起因のエラーを識別する目的で
database層にエラーのインターフェースを定義しています。
```
type ExternalError interface {
	Code() ErrorCode
	Messages() []string
}

type InternalError interface {
	Internal() bool
}
```
新たにrepository層を定義する場合には、
interfaceを実装した構造体を使用して以下のようなイメージで使用してください
```
ExternalDBError{
		ErrorMessage:  fmt.Sprintf(`該当の募集(ID=%d)は見つかりません`, b.ID),
		OriginalError: err,
		StatusCode:    ValueNotFound,
	}
```


# データベースについて
このAPIにおいてデータベースは3つ登場します。

１つ目がherokuにおけるDB、２つ目がローカルにおいて作動するgolang_mysql_dev、そして３つ目がテスト用に用意されたgolang_mysql_testです。

テスト用のデータベースはgolangにおいてテストDBのmock化を調査したところ、handlerなどを含むテストに適切なライブラリが存在しないことが判明したため用意しました。

なおテストDBのcreateはdocker-compose up dbの初回実行時に行われ、その後テーブル内のデータはテスト実行時に初期化されます。

# OpenAPIについて
OpenAPIにはswaggerを採用しました。

またgolangにおいてはswagというコメントに特定の書式を入力することでswagger.ymlを自動記述するライブラリが存在したため、こちらも今回は使用しています。

# herokuデプロイとコンテナ化について
今回は素早くデプロイが出来て、なおかつ他の方のマシンでも即座に動作するものを作りたいと考えていたので、サーバーにはHeroku、コンテナ化にはdocker/docker-composeを採用しました。
