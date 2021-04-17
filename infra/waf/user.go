package waf

import (
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/SakagamiKazuto/golang_api/infra/waf/apperror"
	"github.com/SakagamiKazuto/golang_api/interface/controller"
	"github.com/labstack/echo/v4"
	"net/http"
)

// SignUp is creating user.
// @Summary create user
// @Description create user in a group
// @Accept  json
// @Produce  json
// @Success 201 {object} domain.User
// @Failure 400,409 {object} echo.HTTPError
// @Router /signup [post]
func (s Server) SignUp(uc *controller.UserController) echo.HandlerFunc {
	return C(func(c *Context) error {
		user := new(domain.User)
		if err := c.Bind(user); err != nil {
			return apperror.ResponseError(c, err)
		}

		if err := ValidateUser(user); err != nil {
			return apperror.ResponseError(c, err)
		}

		u, err := uc.CreateUser(user)
		if err != nil {
			return apperror.ResponseError(c, err)
		}
		return c.JSON(http.StatusCreated, u)
	})

}

// Login is creating jwt token.
// @Summary create jwt token
// @Description create token in a group
// @Accept  json
// @Produce  json
// @success 200 {object} object{token=string} "jwt token which you can use to request /api routings"
// @Failure 400,401 {object} echo.HTTPError
// @Router /login [post]
func (s Server) Login(uc *controller.UserController) echo.HandlerFunc {
	return C(func(c *Context) error {
		u := new(domain.User)
		if err := c.Bind(u); err != nil {
			return apperror.ResponseError(c, err)
		}

		user, err := uc.GetUserByMailPass(&domain.User{Password: u.Password, Mail: u.Mail})
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
	})
}
