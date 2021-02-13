package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"work/db"
	"work/model"
)

type jwtCustomClaims struct {
	UID  uint    `json:"uid"`
	Mail string `json:"mail"`
	jwt.StandardClaims
}

//ソースコード公開用
//デプロイ時は"secret"を書き換える
var signingKey = []byte("secret")
//var signingKey = []byte(os.Getenv("SIGNINGKEY"))

var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

// Signup is creating user.
// @Summary create user
// @Description create user in a group
// @Accept  json
// @Produce  json
// @Success 201 {object} model.User
// @Failure 400,409 {object} echo.HTTPError
// @Router /signup [post]
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

// Login is creating jwt token.
// @Summary create jwt token
// @Description create token in a group
// @Accept  json
// @Produce  json
// @success 200 {body} string "jwt token which you can use to request /api routings"
// @Failure 400,401 {object} echo.HTTPError
// @Router /signup [post]
func Login(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	user := model.FindUser(&model.User{Mail: u.Mail, Password: u.Password}, db.DB)
	if user.ID == 0 {
		return &echo.HTTPError{
			Code: http.StatusUnauthorized,
			Message: fmt.Sprintf("can't find user.\n Mail:%v, Pass:%v", u.Mail, u.Password),
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

// リクエストのBodyからTokenを受け取ってUser_idに変換して返す
func userIDFromToken(c echo.Context) uint {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uid := claims.UID
	return uid
}
