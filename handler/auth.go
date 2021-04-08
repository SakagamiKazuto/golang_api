package handler

import (
	"github.com/SakagamiKazuto/golang_api/apperror"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/SakagamiKazuto/golang_api/db"
	"github.com/SakagamiKazuto/golang_api/model"
)

type jwtCustomClaims struct {
	UID  uint   `json:"uid"`
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
		return apperror.ResponseError(c, err)
	}

	if err := user.Validate(); err != nil {
		return apperror.ResponseError(c, err)
	}

	u, err := model.CreateUser(user, db.DB)
	if err != nil {
		return apperror.ResponseError(c, err)
	}
	return c.JSON(http.StatusCreated, u)
}

// Login is creating jwt token.
// @Summary create jwt token
// @Description create token in a group
// @Accept  json
// @Produce  json
// @success 200 {object} object{token=string} "jwt token which you can use to request /api routings"
// @Failure 400,401 {object} echo.HTTPError
// @Router /login [post]
func Login(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return apperror.ResponseError(c, err)
	}

	user, err := model.FindUser(u, db.DB)
	if err != nil {
		return apperror.ResponseError(c, err)
	}

	t, err := CreateToken(user.ID, user.Mail)
	if err != nil {
		return apperror.ResponseError(c, err)
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

func isLogined(c echo.Context) (bool, error) {
	uid := userIDFromToken(c)
	if user, err := model.FindUser(&model.User{Model: gorm.Model{ID: uid}}, db.DB); user.ID == 0 {
		return false, err
	}
	return true, nil
}
