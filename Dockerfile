FROM golang:latest

RUN mkdir /go/src/work

WORKDIR /go/src/work
ADD ./ /go/src/work

RUN go get -u github.com/labstack/echo \
&& go get -u github.com/labstack/echo/middleware \
        && go get github.com/jinzhu/gorm \
        && go get github.com/go-sql-driver/mysql \
        && go get github.com/stretchr/testify/assert \ 
        && go get github.com/dgrijalva/jwt-go \
        && go get github.com/DATA-DOG/go-sqlmock \
        && go get -u github.com/swaggo/swag/cmd/swag \
        && go get github.com/alecthomas/template \
        && go get -u github.com/swaggo/echo-swagger

ENV PORT=${PORT}

CMD ["go", "run", "main.go"]
