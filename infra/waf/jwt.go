package waf

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

type jwtCustomClaims struct {
	UID  uint   `json:"uid"`
	Mail string `json:"mail"`
	jwt.StandardClaims
}

//ソースコード公開用
//デプロイ時は"secret"を書き換える
var signingKey = []byte("secret")

var JwtConf = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

func CreateToken(ID uint, Mail string) (string, error) {
	claims := &jwtCustomClaims{
		ID,
		Mail,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	return t, err
}

//リクエストのBodyからTokenを受け取ってUser_idに変換して返す
func userIDFromToken(c echo.Context) uint {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uid := claims.UID
	return uid
}
