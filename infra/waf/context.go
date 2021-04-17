package waf

import "github.com/labstack/echo/v4"

type Context struct {
	echo.Context
}

type callFunc func(c *Context) error

func C(h callFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(c.(*Context))
	}
}

