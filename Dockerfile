FROM golang:1.16.2-alpine3.13

RUN mkdir /go/src/work

WORKDIR /go/src/work
ADD ./ /go/src/work

# 特定パッケージのgo getが遅い場合にそのデバッグに時間がかかるため細分化
RUN go get -u github.com/labstack/echo \
&& go get -u github.com/labstack/echo/middleware
RUN go get github.com/jinzhu/gorm
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/stretchr/testify/assert
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/DATA-DOG/go-sqlmock
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go get github.com/alecthomas/template
RUN go get -u github.com/swaggo/echo-swagger
RUN go get github.com/joho/godotenv

ENV PORT=${PORT}

CMD ["go", "run", "main.go"]
