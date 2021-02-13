package handler

import (
	//"log"
	"net/http"
	//"os"
	"time"

	//"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"work/model"
	"work/db"
)

type jwtCustomClaims struct {
	UID  uint    `json:"uid"`
	Mail string `json:"mail"`
	jwt.StandardClaims
}

// "secret"書き換える
var signingKey = []byte("secret")
//var signingKey = []byte(os.Getenv("SIGNINGKEY"))

var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

/*
ユーザー登録に際して呼ばれるhandler
既存で存在するか否かを判定後、
DBにuser情報を登録する
*/
func Signup(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	if user.Name == "" || user.Mail == "" || user.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "must write in 'name','email','password'",
		}
	}

	// メールアドレスで重複があれば登録を防ぐ
	if u := model.FindUser(&model.User{Mail: user.Mail}, db.DB); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "this account already exists",
		}
	}

	user = model.CreateUser(user, db.DB)
	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	user := model.FindUser(&model.User{Mail: u.Mail}, db.DB)
	if user.ID == 0 || user.Password != u.Password {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid email or password",
		}
	}

	t, err := CreateToken(user.ID, user.Mail)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
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

// API側で認証の有無を判定するときに必要なロジック
// jwt取れないとき、エラー返すように実装したい
func userIDFromToken(c echo.Context) uint {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uid := claims.UID
	return uid
}
